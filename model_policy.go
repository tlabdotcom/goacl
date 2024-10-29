package goacl

import "github.com/uptrace/bun"

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
