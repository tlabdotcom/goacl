package goacl

import (
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	Name          string    `bun:"name,unique" json:"name"`
	Label         string    `json:"label"`
	Description   string    `json:"description"`
	Features      []Feature `bun:"-" json:"features"`
}

type Feature struct {
	bun.BaseModel `bun:"table:features,alias:f"`
	ID            int64         `bun:"id,pk,autoincrement" json:"id"`
	Name          string        `bun:"name,unique" json:"name"`
	Description   string        `json:"description"`
	SubFeatures   []*SubFeature `bun:"rel:has-many,join:id=feature_id" json:"sub_features,omitempty"`
}

type SubFeature struct {
	bun.BaseModel `bun:"table:sub_features,alias:sf"`
	ID            int64       `bun:"id,pk,autoincrement" json:"id"`
	FeatureID     int64       `bun:"feature_id" json:"feature_id"`
	Feature       *Feature    `bun:"rel:belongs-to,join:feature_id=id" json:"feature,omitempty"`
	Name          string      `bun:"name,unique" json:"name"`
	Description   string      `json:"description"`
	Endpoints     []*Endpoint `bun:"rel:has-many,join:id=sub_feature_id" json:"endpoints,omitempty"`
	Status        bool        `bun:"-" json:"status"`
}

type Endpoint struct {
	bun.BaseModel `bun:"table:endpoints,alias:enp"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	Method        string `bun:"method" json:"method"`
	URL           string `bun:"url" json:"url"`
	SubFeatureID  int64  `bun:"sub_feature_id" json:"sub_feature_id"`
}

type Policy struct {
	bun.BaseModel `bun:"table:policies,alias:p"`
	ID            int64       `bun:"id,pk,autoincrement" json:"id"`
	RoleID        int64       `bun:"role_id" json:"role_id"`
	Role          *Role       `bun:"rel:belongs-to,join:role_id=id"`
	FeatureID     int64       `bun:"feature_id" json:"feature_id"`
	Feature       *Feature    `bun:"rel:belongs-to,join:feature_id=id"`
	SubFeatureID  int64       `bun:"sub_feature_id" json:"sub_feature_id"`
	SubFeature    *SubFeature `bun:"rel:belongs-to,join:sub_feature_id=id" json:"sub_feature,omitempty"`
	Status        bool        `json:"status"`
}

type SubFeatureParam struct {
	ID     int64 `param:"id" query:"id" form:"id" json:"id" xml:"id"`
	Status bool  `param:"status" query:"status" form:"status" json:"status" xml:"status"`
}

type AclParam struct {
	Name        string            `param:"name" query:"name" form:"name" json:"name" xml:"name" validate:"required"`
	Label       string            `param:"label" query:"label" form:"label" json:"label" xml:"label"`
	Description string            `param:"description" query:"description" form:"description" json:"description" xml:"description"`
	SubFeatures []SubFeatureParam `param:"sub_features" query:"sub_features" form:"sub_features" json:"sub_features" xml:"sub_features"`
}

func (p *AclParam) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}
