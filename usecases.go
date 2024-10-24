package goacl

import (
	"context"
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/redis/go-redis/v9"
	"github.com/tlabdotcom/gohelper"
	"github.com/uptrace/bun"
)

const (
	RoleCreatedEvent       string = "AclRoleCreated"
	RoleUpdatedEvent       string = "AclRoleUpdated"
	RoleDeletedEvent       string = "AclRoleDeleted"
	FeatureCreatedEvent    string = "AclFeatureCreated"
	FeatureUpdatedEvent    string = "AclFeatureUpdated"
	FeatureDeletedEvent    string = "AclFeatureDeleted"
	SubFeatureCreatedEvent string = "AclSubfeatureCreated"
	SubFeatureUpdatedEvent string = "AclSubfeatureUpdated"
	SubFeatureDeletedEvent string = "AclSubfeatureDeleted"
	PolicyCreatedEvent     string = "AclPolicyCreated"
	PolicyDeletedEvent     string = "AclPolicyDeleted"
)

type ACL struct {
	DB           *bun.DB
	Redis        *redis.Client
	Enforcer     *casbin.Enforcer
	AclKeyEvents string
}

func NewACL(db *bun.DB, redis *redis.Client, aclKeyEvents *string) (*ACL, error) {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act, dom")
	m.AddDef("g", "g", " _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch5(r.obj,p.obj) && r.act == p.act  || r.sub == 'superadmin'")

	// new adapter
	adapter, err := sqladapter.NewAdapter(db.DB, "postgres", "")
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	if err = enforcer.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
		return nil, err
	}
	if aclKeyEvents == nil {
		aclKeyEvents = gohelper.PointerString("acl_events")
	}
	return &ACL{
		DB:           db,
		Redis:        redis,
		Enforcer:     enforcer,
		AclKeyEvents: gohelper.GetStringValue(aclKeyEvents),
	}, nil
}

// Implement CRUD operations for Feature, SubFeature, and Policy...
func (a *ACL) AddPolicy(ctx context.Context, policy *Policy) error {
	// Publish event
	// a.publishEvent(PolicyCreatedEvent, policy)
	return nil
}
func (a *ACL) RemovePolicy(ctx context.Context, policyID int64) error {
	// policy := &Policy{ID: policyID}
	// a.publishEvent(PolicyDeletedEvent, policy)
	return nil
}

// CRUD operations for Feature
func (a *ACL) CreateFeature(ctx context.Context, feature *Feature) error {
	// a.publishEvent(FeatureCreatedEvent, feature)
	return nil
}
func (a *ACL) UpdateFeature(ctx context.Context, feature *Feature) error {
	// a.publishEvent(FeatureUpdatedEvent, feature)
	return nil
}
func (a *ACL) DeleteFeature(ctx context.Context, featureID int64) error {
	feature := &Feature{ID: featureID}
	_, err := a.DB.NewDelete().Model(feature).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	// a.publishEvent(FeatureDeletedEvent, feature)
	return nil
}

// CRUD operations for SubFeature
func (a *ACL) CreateSubFeature(ctx context.Context, subFeature *SubFeature) error {
	// a.publishEvent(SubFeatureCreatedEvent, subFeature)
	return nil
}
func (a *ACL) UpdateSubFeature(ctx context.Context, subFeature *SubFeature) error {
	// a.publishEvent(SubFeatureUpdatedEvent, subFeature)
	return nil
}
func (a *ACL) DeleteSubFeature(ctx context.Context, subFeatureID int64) error {
	// subFeature := &SubFeature{ID: subFeatureID}
	// a.publishEvent(SubFeatureDeletedEvent, subFeature)
	return nil
}
