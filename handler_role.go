package goacl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tlabdotcom/gohelper"
	"github.com/tlabdotcom/goresponse"
)

// Role handlers
func (a *ACL) listRolesHandler(c echo.Context) error {
	datas, err := a.ListRoles(c.Request().Context())
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(datas, "success", http.StatusOK))
}

func (a *ACL) detailRoleHandler(c echo.Context) error {
	roleID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	data, err := a.GetRoleWithFeatures(c.Request().Context(), roleID)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(data, "success", http.StatusOK))
}

func (a *ACL) createRoleHandler(c echo.Context) error {
	params := new(RoleParam)
	if err := c.Bind(params); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	err := params.Validate()
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	dataData, err := a.CreateRole(c.Request().Context(), params)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}

	data, err := a.GetRoleWithFeatures(c.Request().Context(), dataData.ID)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}

	return c.JSON(http.StatusCreated, goresponse.GenerateSingleDataResponse(data, "Role created successfully", http.StatusCreated))
}

func (a *ACL) updateRoleHandler(c echo.Context) error {
	role := new(RoleParam)
	if err := c.Bind(role); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	err := a.UpdateRole(c.Request().Context(), role)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	data, err := a.GetRoleWithFeatures(c.Request().Context(), role.ID)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(data, "Role updated successfully", http.StatusOK))
}

func (a *ACL) deleteRoleHandler(c echo.Context) error {
	roleID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.DeleteRole(c.Request().Context(), roleID); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(fmt.Sprintf("Role %d deleted", roleID), "Role deleted successfully", http.StatusOK))
}
