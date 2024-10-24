package goacl

import "context"

func (a *ACL) GetSubFeatureByID(ctx context.Context, id int64) (*SubFeature, error) {
	data := &SubFeature{ID: id}
	err := a.DB.NewSelect().Model(data).Relation("Endpoints").WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
