package goacl

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func (a *ACL) GetSubFeatureIDsByFeatureID(ctx context.Context, featureID int64) ([]int64, error) {
	var subFeatureIDs []int64
	err := a.DB.NewSelect().
		Model((*SubFeature)(nil)).
		Column("id").
		Where("feature_id = ?", featureID).
		Scan(ctx, &subFeatureIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sub feature IDs for feature_id %d: %w", featureID, err)
	}
	return subFeatureIDs, nil
}

func (a *ACL) ListSubFeatures(ctx context.Context) ([]*SubFeature, error) {
	var datas []*SubFeature
	err := a.DB.NewSelect().Model(&datas).Relation("Endpoints").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return datas, err
}

func (a *ACL) GetSubFeatureByID(ctx context.Context, id int64) (*SubFeature, error) {
	data := &SubFeature{ID: id}
	err := a.DB.NewSelect().
		Model(data).
		Relation("Endpoints").
		WherePK().
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SubFeature
func (a *ACL) getSubFeatureIncludeEndpointsByIDs(ctx context.Context, ids []int64) ([]SubFeature, error) {
	subFeatures := []SubFeature{}
	err := a.DB.NewSelect().
		Model(&subFeatures).
		Relation("Feature").
		Relation("Endpoints").
		Where("sf.id IN (?)", bun.In(ids)).
		OrderExpr("sf.id ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return subFeatures, nil
}

// func (a *ACL) CreateSubFeature(ctx context.Context, subFeature *SubFeature) error {
// 	_, err := a.DB.NewInsert().Model(subFeature).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (a *ACL) UpdateSubFeature(ctx context.Context, subFeature *SubFeature) error {
	_, err := a.DB.NewUpdate().Model(subFeature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) DeleteSubFeature(ctx context.Context, id int64) error {
	subFeature := &SubFeature{ID: id}
	err := a.DB.NewSelect().Model(subFeature).WherePK().Scan(ctx)
	if err != nil {
		return err
	}
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.NewDelete().Model((*Endpoint)(nil)).Where("sub_feature_id=?", id).Exec(ctx)
	if err != nil {
		return err
	}

	// delete in policies
	_, err = tx.NewDelete().Model((*Policy)(nil)).Where("feature_id=?", subFeature.FeatureID).Where("sub_feature_id=?", id).Exec(ctx)
	if err != nil {
		return err
	}

	// delete data in casbin_rules
	err = a.DeleteCasbinRuleByFeatureAndSubFeatureID(ctx, tx, subFeature.FeatureID, id)
	if err != nil {
		return err
	}

	// delete sub feature
	_, err = tx.NewDelete().Model(subFeature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	// load new policies
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}
