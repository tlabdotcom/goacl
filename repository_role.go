package goacl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

func (a *ACL) CreateRole(ctx context.Context, param *RoleParam) (*Role, error) {
	// Start a transaction
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Step 1: Add the role inside the transaction
	role := &Role{
		Name:        param.Name,
		Label:       param.Label,
		Description: param.Description,
	}
	_, err = tx.NewInsert().Model(role).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to insert role: %w", err)
	}

	// Step 2: Fetch the sub-features with the corresponding IDs
	subIDs := make([]int64, len(param.SubFeatures))
	for i, sub := range param.SubFeatures {
		subIDs[i] = sub.ID
	}

	subs, err := a.getSubFeatureIncludeEndpointsByIDs(ctx, subIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sub-features: %w", err)
	}

	// Step 3: Create policies
	var policies []Policy
	for _, subFeature := range subs {
		for range subFeature.Endpoints {
			policy := Policy{
				RoleID:       role.ID,
				FeatureID:    subFeature.FeatureID,
				SubFeatureID: subFeature.ID,
				Status:       true,
			}
			policies = append(policies, policy)
		}
	}

	// Step 4: Insert policies into the database inside the transaction
	// Capture the auto-incremented IDs
	_, err = tx.NewInsert().Model(&policies).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to insert policies: %w", err)
	}

	// Step 5: Create Casbin rules with policy IDs
	var policiesRules [][]string
	for i, subFeature := range subs {
		for _, endpoint := range subFeature.Endpoints {
			policiesRules = append(policiesRules, []string{
				role.Name,
				endpoint.URL,
				endpoint.Method,
				fmt.Sprint(policies[i].ID),
			})
		}
	}

	// Step 6: Add Casbin policies (using Casbin's Enforcer)
	ok, err := a.Enforcer.AddPolicies(policiesRules)
	if err != nil {
		return nil, fmt.Errorf("failed to add Casbin policies: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("some Casbin policies were not added")
	}
	// Step 7: Commit the transaction (handled by defer)
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	fmt.Println("Policies and rules added successfully:", ok)
	return role, nil
}

func (a *ACL) ListRoles(ctx context.Context) ([]Role, error) {
	datas := []Role{}
	err := a.DB.NewSelect().Model(&datas).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (a *ACL) UpdateRole(ctx context.Context, param *RoleParam) error {
	// Start a transaction
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Fetch the existing role data
	roleData, err := a.GetRoleByID(ctx, param.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch role: %w", err)
	}

	// Validate the role for update
	updatedRole, err := param.ValidateForUpdate(roleData)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update the role details
	_, err = tx.NewUpdate().Model(updatedRole).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Step 2: Process the sub-features and update policies
	var policiesToAdd []Policy
	var policiesToRemove []string
	var policiesRulesToAdd [][]string
	var policiesRulesToRemove [][]string

	// fmt.Println("data: ", subs)

	for _, subFeature := range param.SubFeatures {
		newSub, err := a.GetSubFeatureByID(ctx, subFeature.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		for _, endpoint := range newSub.Endpoints {
			if subFeature.Status {
				// Add policies when the sub-feature status is true
				policy := Policy{
					RoleID:       updatedRole.ID,
					FeatureID:    newSub.FeatureID,
					SubFeatureID: subFeature.ID,
					Status:       true,
				}
				policiesToAdd = append(policiesToAdd, policy)

				// Casbin rule to add (with policy ID placeholder, to be updated after insert)
				policiesRulesToAdd = append(policiesRulesToAdd, []string{
					updatedRole.Name,
					endpoint.URL,
					endpoint.Method,
					fmt.Sprint(policy.ID), // Policy ID will be updated after insert
				})
			} else {
				// Fetch existing policies for the sub-feature
				existingPolicies, err := a.GetPoliciesByRoleAndSubFeature(ctx, updatedRole.ID, subFeature.ID)
				if err != nil {
					return fmt.Errorf("failed to fetch policies: %w", err)
				}

				// Collect policies to remove
				for _, policy := range existingPolicies {
					policiesToRemove = append(policiesToRemove, fmt.Sprint(policy.ID))
					// Casbin rule to remove (with policy ID)
					policiesRulesToRemove = append(policiesRulesToRemove, []string{
						updatedRole.Name,
						endpoint.URL,
						endpoint.Method,
						fmt.Sprint(policy.ID),
					})
				}
			}
		}
	}

	// Step 3: Remove policies and Casbin rules for disabled sub-features
	if len(policiesToRemove) > 0 {
		_, err = tx.NewDelete().Model((*Policy)(nil)).
			Where("id IN (?)", bun.In(policiesToRemove)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete policies: %w", err)
		}

		// Remove Casbin rules
		_, err = a.Enforcer.RemovePolicies(policiesRulesToRemove)
		if err != nil {
			return fmt.Errorf("failed to remove Casbin policies: %w", err)
		}
	}

	// Step 4: Add new policies for enabled sub-features
	if len(policiesToAdd) > 0 {
		// Insert policies into the database
		_, err = tx.NewInsert().Model(&policiesToAdd).Returning("id").Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert policies: %w", err)
		}

		// Update the policy IDs in the Casbin rules to add
		for i, policy := range policiesToAdd {
			// Update the last element (policy ID) of the corresponding rule
			policiesRulesToAdd[i][3] = fmt.Sprint(policy.ID)
		}

		// Add Casbin rules
		_, err = a.Enforcer.AddPolicies(policiesRulesToAdd)
		if err != nil {
			return fmt.Errorf("failed to add Casbin policies: %w", err)
		}
	}

	// Step 5: Commit the transaction (handled by defer)
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to reload Casbin policies: %w", err)
	}

	return nil
}

func (a *ACL) GetRoleWithFeatures(ctx context.Context, roleID int64) (*Role, error) {
	role := &Role{ID: roleID}
	// Query to fetch Role with Features and SubFeatures
	err := a.DB.NewSelect().
		Model(role).
		WherePK().
		Scan(ctx)
	// If no data is found, return a graceful message
	if err != nil {
		return nil, err
	}

	features := []Feature{}
	err = a.DB.NewSelect().Model(&features).Relation("SubFeatures").Scan(ctx)
	if err != nil {
		return nil, err
	}
	role.Features = features
	// Fetch the policy status for each sub-feature

	policies := []Policy{}
	_, err = a.DB.NewSelect().
		Model(&policies).
		Where("role_id = ?", roleID).
		ScanAndCount(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	fmt.Println(policies)
	// Map the policy status to sub-features
	statusMap := make(map[int64]bool)
	for _, policy := range policies {
		statusMap[policy.SubFeatureID] = policy.Status
	}

	// Assign status to each sub-feature
	for i, feature := range role.Features {
		for j, subFeature := range feature.SubFeatures {
			subFeature.Status = statusMap[subFeature.ID]
			role.Features[i].SubFeatures[j] = subFeature
		}
	}
	return role, nil
}

func (a *ACL) DeleteRole(ctx context.Context, roleID int64) error {
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	role := &Role{ID: roleID}
	_, err = tx.NewDelete().Model(role).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	policyIDs, err := a.GetPolicyIDsByRoleID(ctx, roleID)
	if err != nil {
		return err
	}
	// delete policy to Casbin by Role
	stringPolicyIDs := make([]string, len(policyIDs))
	for i, id := range policyIDs {
		stringPolicyIDs[i] = fmt.Sprintf("%d", id)
	}

	_, err = tx.NewDelete().
		TableExpr("casbin_rule").
		Where("v3 IN (?)", bun.In(stringPolicyIDs)).
		Exec(ctx)
	if err != nil {
		fmt.Println("Error while deleting records:", err)
		return err
	}

	// delete policy to DB by Role
	_, err = tx.NewDelete().Model(&Policy{}).Where("role_id=?", roleID).Exec(ctx)
	if err != nil {
		return err
	}
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	fmt.Println("Policies and rules deleted successfully:")
	return nil
}

func (a *ACL) GetRoleByID(ctx context.Context, roleID int64) (*Role, error) {
	role := &Role{ID: roleID}
	// Query to fetch Role with Features and SubFeatures
	err := a.DB.NewSelect().
		Model(role).
		WherePK().
		Scan(ctx)
	// If no data is found, return a graceful message
	fmt.Println(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
