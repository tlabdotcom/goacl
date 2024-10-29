package goacl

import (
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`
	ID            int64      `bun:"id,pk,autoincrement" json:"id"`
	Name          string     `bun:"name,unique" json:"name,omitempty"`
	Label         string     `json:"label,omitempty"`
	Description   string     `json:"description,omitempty"`
	Features      []*Feature `bun:"-" json:"features,omitempty"`
}

type RoleParam struct {
	ID          int64             `param:"id" query:"id" form:"id" json:"id" xml:"id"`
	Name        string            `param:"name" query:"name" form:"name" json:"name" xml:"name" validate:"required"`
	Label       string            `param:"label" query:"label" form:"label" json:"label" xml:"label" validate:"required"`
	Description string            `param:"description" query:"description" form:"description" json:"description" xml:"description"`
	SubFeatures []SubFeatureParam `param:"sub_features" query:"sub_features" form:"sub_features" json:"sub_features" xml:"sub_features"`
}

func (p *RoleParam) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *RoleParam) ValidateForUpdate(data *Role) (*Role, error) {
	if p.Label != "" {
		data.Label = p.Label
	}
	if p.Description != "" {
		data.Description = p.Description
	}
	return data, nil
}
