package goacl

import (
	"context"
	"log"
)

func seedRoles(ctx context.Context, a *ACL, roles []Role) error {
	_, err := a.DB.NewInsert().Model(&roles).On("CONFLICT (name) DO NOTHING").Exec(ctx)
	if err != nil {
		log.Printf("Error seeding role: %v", err)
		return err
	}
	return nil
}

func seedFeatures(ctx context.Context, a *ACL, features []Feature) error {
	_, err := a.DB.NewInsert().Model(&features).On("CONFLICT (name) DO NOTHING").Exec(ctx)
	if err != nil {
		log.Printf("Error seeding feature: %v", err)
		return err
	}

	return nil
}

func seedSubFeatures(ctx context.Context, a *ACL, subFeatures []SubFeature) error {
	_, err := a.DB.NewInsert().Model(&subFeatures).On("CONFLICT (name) DO NOTHING").Exec(ctx)
	if err != nil {
		log.Printf("Error seeding sub-feature: %v", err)
		return err
	}
	return nil
}

func seedEndpoints(ctx context.Context, a *ACL, endpoints []Endpoint) error {
	_, err := a.DB.NewInsert().Model(&endpoints).On("CONFLICT DO NOTHING").Exec(ctx)
	if err != nil {
		log.Printf("Error seeding endpoint: %v", err)
		return err
	}

	return nil
}

func seedPolicies(ctx context.Context, a *ACL, policies []Policy) error {
	_, err := a.DB.NewInsert().Model(&policies).On("CONFLICT DO NOTHING").Exec(ctx)
	if err != nil {
		log.Printf("Error seeding policy: %v", err)
		return err
	}
	return nil
}
