package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/tlabdotcom/goacl"
	"github.com/tlabdotcom/godb"
)

func init() {
	viper.SetConfigName("../test_env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func main() {
	e := echo.New()
	// initial go acl
	acl, err := goacl.NewACL(godb.GetPostgresDB(), godb.GetRedis(), nil)
	if err != nil {
		e.Logger.Fatal(err)
	}
	// consummer listener
	go acl.ConsumeEvents(context.TODO())
	// setup routes
	acl.SetupRoutes(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
