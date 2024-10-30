package goacl

import (
	"github.com/labstack/echo/v4"
)

func (a *ACL) SetupRoutes(e *echo.Echo) {
	g := e.Group("/acl")
	g.GET("/role/list", a.listRolesHandler)
	g.POST("/role/create", a.createRoleHandler)
	g.GET("/role/detail/:id", a.detailRoleHandler)
	g.PUT("/role/update/:id", a.updateRoleHandler)
	g.DELETE("/role/delete/:id", a.deleteRoleHandler)

	g.GET("/feature/list", a.listFeatureHandler)
	g.POST("/feature/create", a.createFeatureHandler)
	g.GET("/feature/detail/:id", a.detailFeatureHandler)
	// g.PUT("/feature/updated/:id", a.updateFeatureHandler)
	g.DELETE("/feature/delete/:id", a.deleteFeatureHandler)

	// g.GET("/sub-feature/list", a.listSubFeatureHandler)
	// g.POST("/sub-feature/create", a.createSubFeatureHandler)
	g.PUT("/subfeature/update/:id", a.updateSubFeatureHandler)
	g.DELETE("/subfeature/delete/:id", a.deleteSubFeatureHandler)

	// enpoints
	g.DELETE("/endpoint/delete/:id", a.deleteEndpointHandler)
}
