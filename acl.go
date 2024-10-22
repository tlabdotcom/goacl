package goacl

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/uptrace/bun"
)

// Similar handler functions for update and delete...
func (a *ACL) CheckPermission(ctx context.Context, roleName, url, method string) (bool, error) {
	return a.Enforcer.Enforce(roleName, url, method)
}

// Helper function to load policies from the database
func loadPoliciesFromDB(enforcer *casbin.Enforcer, db *bun.DB) error {
	var policies []Policy
	err := db.NewSelect().Model(&policies).Relation("Role").Relation("Feature").Relation("SubFeature", func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.Relation("Endpoint")
	}).Scan(context.Background())
	if err != nil {
		return err
	}
	// TODO
	for _, policy := range policies {
		for _, item := range policy.SubFeature.Endpoints {
			_, err := enforcer.AddPolicy(policy.Role.Name, policy.Feature.Name, item.URL, item.Method)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
