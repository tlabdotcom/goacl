package goacl

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func (a *ACL) DeleteCasbinRuleByFeatureAndSubFeatureID(ctx context.Context, tx bun.Tx, featureID, subFeatureID int64) error {
	policyIDs, err := a.GetPolicyIDsByFeatureIDAndSubFeatureID(ctx, featureID, subFeatureID)
	if err != nil {
		return err
	}
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
	return nil
}
