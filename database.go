package goacl

import "context"

func (a *ACL) Migrate(ctx context.Context) error {
	models := []interface{}{
		(*Role)(nil),
		(*Feature)(nil),
		(*SubFeature)(nil),
		(*Endpoint)(nil),
		(*Policy)(nil),
	}

	for _, model := range models {
		_, err := a.DB.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ACL) Seed(
	ctx context.Context,
	roles []Role,
	features []Feature,
	subFeatures []SubFeature,
	endpoints []Endpoint,
	policies []Policy,
) error {
	// Insert the sample data into the database
	if err := seedRoles(ctx, a, roles); err != nil {
		return err
	}
	if err := seedFeatures(ctx, a, features); err != nil {
		return err
	}
	if err := seedSubFeatures(ctx, a, subFeatures); err != nil {
		return err
	}
	if err := seedEndpoints(ctx, a, endpoints); err != nil {
		return err
	}
	if err := seedPolicies(ctx, a, policies); err != nil {
		return err
	}
	return nil
}
