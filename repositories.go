package goacl

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
)

func (a *ACL) createRoleToDB(ctx context.Context, param *AclParam) error {
	role := &Role{
		Name:        param.Name,
		Label:       param.Label,
		Description: param.Description,
	}
	_, err := a.DB.NewInsert().Model(role).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) updateRoleToDB(ctx context.Context, role *Role) error {
	_, err := a.DB.NewUpdate().Model(role).WherePK().Exec(ctx)
	if err != nil {
		log.Errorf("Error updating role: %v", err)
	}
	return nil
}

// Feature
func (a *ACL) listFeatures(ctx context.Context) ([]Feature, error) {
	features := []Feature{}
	err := a.DB.NewSelect().Model(&features).Relation("SubFeatures").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return features, nil
}
func (a *ACL) createFeatureToDB(ctx context.Context, feature *Feature) error {
	_, err := a.DB.NewInsert().Model(feature).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) updateFeatureToDB(ctx context.Context, feature *Feature) error {
	_, err := a.DB.NewUpdate().Model(feature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) deleteFeatureToDB(ctx context.Context, feature *Feature) error {
	_, err := a.DB.NewDelete().Model(feature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
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

func (a *ACL) createSubFeatureToDB(ctx context.Context, subFeature *SubFeature) error {
	_, err := a.DB.NewInsert().Model(subFeature).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) updateSubFeatureToDB(ctx context.Context, subFeature *SubFeature) error {
	_, err := a.DB.NewUpdate().Model(subFeature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *ACL) deleteSubFeatureToDB(ctx context.Context, subFeature *SubFeature) error {
	_, err := a.DB.NewDelete().Model(subFeature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Policy
func (a *ACL) createPolicyToDB(ctx context.Context, policy *Policy) error {
	_, err := a.DB.NewInsert().Model(policy).Exec(ctx)
	if err != nil {
		return err
	}

	// Add policy to Casbin enforcer
	subFeature := &SubFeature{}
	err = a.DB.NewSelect().Model(subFeature).Where("id = ?", policy.SubFeatureID).Scan(ctx)
	if err != nil {
		return err
	}

	role := &Role{}
	err = a.DB.NewSelect().Model(role).Where("id = ?", policy.RoleID).Scan(ctx)
	if err != nil {
		return err
	}

	// _, err = a.Enforcer.AddPolicy(role.Name, subFeature.URL, subFeature.Method)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (a *ACL) deletePolicyToDB(ctx context.Context, policy *Policy) error {
	err := a.DB.NewSelect().Model(policy).WherePK().Scan(ctx)
	if err != nil {
		return err
	}

	_, err = a.DB.NewDelete().Model(policy).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	// Remove policy from Casbin enforcer
	subFeature := &SubFeature{}
	err = a.DB.NewSelect().Model(subFeature).Where("id = ?", policy.SubFeatureID).Scan(ctx)
	if err != nil {
		return err
	}

	role := &Role{}
	err = a.DB.NewSelect().Model(role).Where("id = ?", policy.RoleID).Scan(ctx)
	if err != nil {
		return err
	}

	// _, err = a.Enforcer.RemovePolicy(role.Name, subFeature.URL, subFeature.Method)
	// if err != nil {
	// 	return err
	// }
	return nil
}
