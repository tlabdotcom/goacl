package goacl

import (
	"github.com/uptrace/bun"
)

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

type SubFeatureParam struct {
	ID          int64            `param:"id" query:"id" form:"id" json:"id" xml:"id"`
	Status      bool             `param:"status" query:"status" form:"status" json:"status" xml:"status"`
	FeatureID   int64            `param:"feature_id" query:"feature_id" form:"feature_id" json:"feature_id" xml:"feature_id"`
	Name        string           `param:"name" query:"name" form:"name" json:"name" xml:"name"`
	Description string           `param:"description" query:"description" form:"description" json:"description" xml:"description"`
	Endpoints   []*EndpointParam `param:"endpoints" query:"endpoints" form:"endpoints" json:"endpoints" xml:"endpoints"`
}

func (p *SubFeatureParam) ValidateForUpdate(data *SubFeature) (*SubFeature, error) {
	if p.Name != "" {
		data.Name = p.Name
	}
	if p.Description != "" {
		data.Description = p.Description
	}
	if p.Status {
		data.Status = p.Status
	}
	if p.FeatureID != 0 {
		data.FeatureID = p.FeatureID
	}
	return data, nil
}
