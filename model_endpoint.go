package goacl

import (
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Endpoint struct {
	bun.BaseModel `bun:"table:endpoints,alias:enp"`
	ID            int64  `bun:"id,pk,autoincrement" json:"id"`
	Method        string `bun:"method" json:"method"`
	URL           string `bun:"url" json:"url"`
	SubFeatureID  int64  `bun:"sub_feature_id" json:"sub_feature_id"`
}

type EndpointParam struct {
	ID           int64  `param:"id" query:"id" form:"id" json:"id" xml:"id"`
	Method       string `param:"method" query:"method" form:"method" json:"method" xml:"method"`
	URL          string `param:"url" query:"url" form:"url" json:"url" xml:"url"`
	SubFeatureID int64  `param:"sub_feature_id" query:"sub_feature_id" form:"sub_feature_id" json:"sub_feature_id" xml:"sub_feature_id"`
}

func (p *EndpointParam) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *EndpointParam) ValidateForUpdate(data *Endpoint) (*Endpoint, error) {
	if p.URL != "" {
		data.URL = p.URL
	}
	if p.Method != "" {
		data.Method = p.Method
	}
	return data, nil
}
