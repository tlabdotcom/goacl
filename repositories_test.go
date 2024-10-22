package goacl

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/tlabdotcom/godb"
	"github.com/uptrace/bun"
)

func init() {
	viper.SetConfigName("test_env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func TestACL_GetSubFeatureIncludeEndpointsByIDs(t *testing.T) {
	db := godb.GetPostgresDB()
	rd := godb.GetRedis()
	type fields struct {
		DB           *bun.DB
		Redis        *redis.Client
		AclKeyEvents string
	}
	type args struct {
		ctx context.Context
		ids []int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []SubFeature
		wantErr bool
	}{
		{
			name: "Test Get by IDs",
			fields: fields{
				DB:    db,
				Redis: rd,
			},
			args: args{
				ctx: context.TODO(),
				ids: []int64{1, 2, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewACL(tt.fields.DB, tt.fields.Redis, &tt.fields.AclKeyEvents)
			got, err := a.getSubFeatureIncludeEndpointsByIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("ACL.GetSubFeatureIncludeEndpointsByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
