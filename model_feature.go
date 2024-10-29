package goacl

import (
	"github.com/uptrace/bun"
)

type Feature struct {
	bun.BaseModel `bun:"table:features,alias:f"`
	ID            int64         `bun:"id,pk,autoincrement" json:"id"`
	Name          string        `bun:"name,unique" json:"name"`
	Description   string        `json:"description"`
	SubFeatures   []*SubFeature `bun:"rel:has-many,join:id=feature_id" json:"sub_features,omitempty"`
}

// type FeatureParam struct {
// 	ID          int64         `param:"id" query:"id" form:"id" json:"id" xml:"id"`
// 	Name        string        `param:"name" query:"name" form:"name" json:"name" xml:"name"`
// 	Description string        `param:"description" query:"description" form:"description" json:"description" xml:"description"`
// 	SubFeatures []*SubFeature `param:"sub_features" query:"sub_features" form:"sub_features" json:"sub_features" xml:"sub_features"`
// }

// func (p *FeatureParam) Validate() error {
// 	validate := validator.New()
// 	err := validate.Struct(p)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (p *FeatureParam) ValidateForUpdate(data *Feature) (*Feature, error) {
// 	if p.Name != "" {
// 		data.Name = p.Name
// 	}
// 	if p.Description != "" {
// 		data.Description = p.Description
// 	}
// 	return data, nil
// }
