package main

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/tlabdotcom/goacl"
	"github.com/tlabdotcom/godb"
	"github.com/tlabdotcom/goresponse"
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
	// setup routes
	acl.SetupRoutes(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/inventory", func(c echo.Context) error {
		return c.String(http.StatusOK, "Inventory")
	}, authorizationMiddleware(acl.Enforcer))

	e.GET("/sales", func(c echo.Context) error {
		return c.String(http.StatusOK, "Sales")
	}, authorizationMiddleware(acl.Enforcer))

	e.Logger.Fatal(e.Start(":1323"))
}

func authorizationMiddleware(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Replace with actual user role fetching logic
			role := "user"
			obj := c.Request().URL.Path
			act := c.Request().Method
			// Check permission
			allowed, err := enforcer.Enforce(role, obj, act)
			if err != nil {
				return goresponse.NewStandardErrorResponse(401).AddMessageError("authorization", err.Error()).JSON(c)
			}
			if !allowed {
				return goresponse.NewStandardErrorResponse(403).AddMessageError("authorization", "access denied").JSON(c)
			}

			return next(c)
		}
	}
}
