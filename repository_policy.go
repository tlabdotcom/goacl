package goacl

import (
	"context"
	"fmt"
)

func (a *ACL) GetPolicyIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var policyIDs []int64
	err := a.DB.NewSelect().
		Model((*Policy)(nil)).
		Column("id").
		Where("role_id = ?", roleID).
		Scan(ctx, &policyIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policy IDs for role_id %d: %w", roleID, err)
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
		return nil, fmt.Errorf("failed to fetch policies for role %d and sub-feature %d: %w", roleID, subFeatureID, err)
	}

	return policies, nil
}
