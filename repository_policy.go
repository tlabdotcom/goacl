package goacl

import (
	"context"
)

func (a *ACL) GetPolicyIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var policyIDs []int64
	err := a.DB.NewSelect().
		Model((*Policy)(nil)).
		Column("id").
		Where("role_id = ?", roleID).
		Scan(ctx, &policyIDs)

	if err != nil {
		return nil, err
	}
	return policyIDs, nil
}

func (a *ACL) GetPolicyIDsByFeatureID(ctx context.Context, featureID int64) ([]int64, error) {
	var policyIDs []int64
	err := a.DB.NewSelect().
		Model((*Policy)(nil)).
		Column("id").
		Where("feature_id = ?", featureID).
		Scan(ctx, &policyIDs)

	if err != nil {
		return nil, err
	}
	return policyIDs, nil
}

func (a *ACL) GetPolicyIDsByFeatureIDAndSubFeatureID(ctx context.Context, featureID, subFeatureID int64) ([]int64, error) {
	var policyIDs []int64
	err := a.DB.NewSelect().
		Model((*Policy)(nil)).
		Column("id").
		Where("feature_id = ?", featureID).
		Where("sub_feature_id = ?", subFeatureID).
		Scan(ctx, &policyIDs)

	if err != nil {
		return nil, err
	}
	return policyIDs, nil
}

func (a *ACL) GetPoliciesByRoleAndSubFeature(ctx context.Context, roleID, subFeatureID int64) ([]Policy, error) {
	// Prepare a slice to hold the policies
	var policies []Policy

	// Query the policies for the given role and sub-feature
	err := a.DB.NewSelect().
		Model(&policies).
		Where("role_id = ?", roleID).
		Where("sub_feature_id = ?", subFeatureID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return policies, nil
}
