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

func TestNewACL(t *testing.T) {
	db := godb.GetPostgresDB()
	rd := godb.GetRedis()
	type args struct {
		db           *bun.DB
		redis        *redis.Client
		aclKeyEvents *string
	}
	tests := []struct {
		name    string
		args    args
		want    *ACL
		wantErr bool
	}{
		{
			name: "Test ACL",
			args: args{
				db:    db,
				redis: rd,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := NewACL(tt.args.db, tt.args.redis, tt.args.aclKeyEvents)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewACL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// migrate
			err = got.Migrate(ctx)
			if err != nil {
				t.Error(err)
			}

			// seed
			err = got.Seed(ctx, SampleRoles, SampleFeatures, SampleSubFeatures, SampleEndpoints, SamplePolicies)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
