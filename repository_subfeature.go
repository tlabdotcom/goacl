package goacl

import (
	"context"
	"database/sql"
	"errors"
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

func (a *ACL) UpdateSubFeature(ctx context.Context, params *SubFeatureParam) error {
	subFeature := &SubFeature{ID: params.ID}
	err := a.DB.NewSelect().Model(subFeature).WherePK().Scan(ctx)
	if err != nil {
		return err
	}
	subFeature, err = params.ValidateForUpdate(subFeature)
	if err != nil {
		return err
	}

	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.NewUpdate().Model(subFeature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	if len(params.Endpoints) > 0 {
		for _, endpointParam := range params.Endpoints {
			endpoint := &Endpoint{ID: endpointParam.ID}
			err := tx.NewSelect().Model(endpoint).WherePK().Scan(ctx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}

			updatedEndpoint, err := endpointParam.ValidateForUpdate(endpoint)
			if err != nil {
				return err
			}

			if errors.Is(err, sql.ErrNoRows) || updatedEndpoint.ID == 0 {
				// Check if an endpoint with the same url and sub_feature_id already exists
				existingEndpoint := &Endpoint{}
				err := tx.NewSelect().
					Model(existingEndpoint).
					Where("url = ?", updatedEndpoint.URL).
					Where("sub_feature_id = ?", subFeature.ID).
					Limit(1).
					Scan(ctx)

				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return err
				}

				if existingEndpoint.ID != 0 {
					// Update the existing endpoint if a duplicate is found
					updatedEndpoint.ID = existingEndpoint.ID
					updatedEndpoint.SubFeatureID = subFeature.ID
					if _, err = tx.NewUpdate().Model(updatedEndpoint).WherePK().Exec(ctx); err != nil {
						return err
					}
				} else {
					// Insert if no duplicate is found
					updatedEndpoint.SubFeatureID = subFeature.ID
					if _, err := tx.NewInsert().Model(updatedEndpoint).Exec(ctx); err != nil {
						return err
					}
				}
			} else {
				// Update existing endpoint
				if _, err = tx.NewUpdate().Model(updatedEndpoint).WherePK().Exec(ctx); err != nil {
					return err
				}
			}
		}
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
