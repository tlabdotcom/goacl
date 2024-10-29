package goacl

import (
	"context"
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type ACL struct {
	DB       *bun.DB
	Redis    *redis.Client
	Enforcer *casbin.Enforcer
}

func NewACL(db *bun.DB, redis *redis.Client, aclKeyEvents *string) (*ACL, error) {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act, dom")
	m.AddDef("g", "g", " _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch5(r.obj,p.obj) && r.act == p.act  || r.sub == 'superadmin'")

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
	newACL := &ACL{
		DB:       db,
		Redis:    redis,
		Enforcer: enforcer,
	}
	err = newACL.Migrate(context.TODO())
	if err != nil {
		return nil, err
	}
	return newACL, nil
}

// Similar handler functions for update and delete...
func (a *ACL) CheckPermission(ctx context.Context, roleName, url, method string) (bool, error) {
	return a.Enforcer.Enforce(roleName, url, method)
}
