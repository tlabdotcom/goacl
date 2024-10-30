package goacl

import (
	"context"
)

func (a *ACL) DeleteEndpoint(ctx context.Context, id int64) error {
	endpoint := &Endpoint{ID: id}
	err := a.DB.NewSelect().
		Model(endpoint).
		Relation("SubFeature").
		WherePK().
		Scan(ctx)
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

	err = a.DeleteCasbinRuleByFeatureAndSubFeatureID(ctx, tx, endpoint.SubFeature.FeatureID, endpoint.SubFeatureID)
	if err != nil {
		return err
	}

	// delete endpoints by id
	_, err = tx.NewDelete().Model(endpoint).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
