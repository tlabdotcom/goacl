package goacl

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

// CRUD operations for Feature
func (a *ACL) CreateFeature(ctx context.Context, feature *Feature) error {
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

	_, err = tx.NewInsert().Model(feature).Exec(ctx)
	if err != nil {
		return err
	}

	for _, subFeature := range feature.SubFeatures {
		subFeature.FeatureID = feature.ID
		_, err = tx.NewInsert().Model(subFeature).Exec(ctx)
		if err != nil {
			return err
		}

		for _, endp := range subFeature.Endpoints {
			endp.SubFeatureID = subFeature.ID
			_, err = tx.NewInsert().Model(endp).Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ListFeatures
func (a *ACL) ListFeatures(ctx context.Context) ([]Feature, error) {
	features := []Feature{}
	err := a.DB.NewSelect().
		Model(&features).
		Relation("SubFeatures").
		Relation("SubFeatures.Endpoints").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return features, nil
}

func (a *ACL) DetailFeature(ctx context.Context, id int64) (*Feature, error) {
	feature := &Feature{ID: id}
	err := a.DB.NewSelect().
		Model(feature).
		Relation("SubFeatures").
		Relation("SubFeatures.Endpoints").
		WherePK().
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

// func (a *ACL) UpdateFeature(ctx context.Context, feature *Feature) error {
// 	tx, err := a.DB.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 			return
// 		}
// 		err = tx.Commit()
// 	}()

// 	_, err = tx.NewUpdate().Model(feature).WherePK().Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	for _, subFeature := range feature.SubFeatures {
// 		if subFeature.ID != 0 {
// 			_, err = tx.NewUpdate().Model(subFeature).WherePK().Exec(ctx)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			subFeature.FeatureID = feature.ID
// 			_, err = tx.NewInsert().Model(subFeature).Exec(ctx)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

func (a *ACL) DeleteFeature(ctx context.Context, featureID int64) error {
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

	subFeatureIDs, err := a.GetSubFeatureIDsByFeatureID(ctx, featureID)
	if err != nil {
		return err
	}
	// delete endpoint data by sub feature id
	stringSubFeatureIDs := make([]string, len(subFeatureIDs))
	for i, id := range subFeatureIDs {
		stringSubFeatureIDs[i] = fmt.Sprintf("%d", id)
	}
	_, err = tx.NewDelete().
		Model((*Endpoint)(nil)).
		Where("sub_feature_id IN (?)", bun.In(stringSubFeatureIDs)).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = tx.NewDelete().Model(&SubFeature{}).Where("feature_id = ?", featureID).Exec(ctx)
	if err != nil {
		return err
	}

	feature := &Feature{ID: featureID}
	_, err = tx.NewDelete().Model(feature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	// get policy ids
	policyIDs, err := a.GetPolicyIDsByFeatureID(ctx, featureID)
	if err != nil {
		return err
	}

	// delete policy to Casbin by Role
	stringPolicyIDs := make([]string, len(policyIDs))
	for i, id := range policyIDs {
		stringPolicyIDs[i] = fmt.Sprintf("%d", id)
	}

	_, err = tx.NewDelete().
		TableExpr("casbin_rule").
		Where("v3 IN (?)", bun.In(stringPolicyIDs)).
		Exec(ctx)
	if err != nil {
		return err
	}
	// delete policy to DB by Feature
	_, err = tx.NewDelete().Model(&Policy{}).Where("feature_id=?", featureID).Exec(ctx)
	if err != nil {
		return err
	}
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}
