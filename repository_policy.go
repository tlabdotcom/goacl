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
