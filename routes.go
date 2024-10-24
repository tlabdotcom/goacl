package goacl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tlabdotcom/gohelper"
	"github.com/tlabdotcom/goresponse"
)

func (a *ACL) SetupRoutes(e *echo.Echo) {
	g := e.Group("/acl")
	g.GET("/roles", a.listRolesHandler)
	g.POST("/roles", a.createRoleHandler)
	g.GET("/roles/:id", a.detailRoleHandler)
	g.PUT("/roles/:id", a.updateRoleHandler)
	g.DELETE("/roles/:id", a.deleteRoleHandler)

	g.GET("/features", a.listFeatureHandler)
	g.POST("/features", a.createFeatureHandler)
	g.PUT("/features/:id", a.updateFeatureHandler)
	g.DELETE("/features/:id", a.deleteFeatureHandler)

	g.POST("/subfeatures", a.createSubFeatureHandler)
	g.PUT("/subfeatures/:id", a.updateSubFeatureHandler)
	g.DELETE("/subfeatures/:id", a.deleteSubFeatureHandler)

	g.POST("/policies", a.addPolicyHandler)
	g.DELETE("/policies/:id", a.removePolicyHandler)
}

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
	data, err := a.CreateRole(c.Request().Context(), params)
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
	var err error
	role.ID, err = gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	err = a.UpdateRole(c.Request().Context(), role)
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

// Feature handlers
func (a *ACL) listFeatureHandler(c echo.Context) error {
	datas, err := a.listFeatures(c.Request().Context())
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(datas, "success", http.StatusOK))
}

func (a *ACL) createFeatureHandler(c echo.Context) error {
	feature := new(Feature)
	if err := c.Bind(feature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.CreateFeature(c.Request().Context(), feature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusCreated, goresponse.GenerateSingleDataResponse(feature, "Feature created successfully", http.StatusCreated))
}

func (a *ACL) updateFeatureHandler(c echo.Context) error {
	feature := new(Feature)
	if err := c.Bind(feature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	var err error
	feature.ID, err = gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.UpdateFeature(c.Request().Context(), feature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(feature, "Feature updated successfully", http.StatusOK))
}

func (a *ACL) deleteFeatureHandler(c echo.Context) error {
	featureID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.DeleteFeature(c.Request().Context(), featureID); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(fmt.Sprintf("Feature %d deleted", featureID), "Feature deleted successfully", http.StatusOK))
}

// SubFeature handlers
func (a *ACL) createSubFeatureHandler(c echo.Context) error {
	subFeature := new(SubFeature)
	if err := c.Bind(subFeature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.CreateSubFeature(c.Request().Context(), subFeature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusCreated, goresponse.GenerateSingleDataResponse(subFeature, "SubFeature created successfully", http.StatusCreated))
}

func (a *ACL) updateSubFeatureHandler(c echo.Context) error {
	subFeature := new(SubFeature)
	if err := c.Bind(subFeature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	var err error
	subFeature.ID, err = gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.UpdateSubFeature(c.Request().Context(), subFeature); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(subFeature, "SubFeature updated successfully", http.StatusOK))
}

func (a *ACL) deleteSubFeatureHandler(c echo.Context) error {
	subFeatureID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.DeleteSubFeature(c.Request().Context(), subFeatureID); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(fmt.Sprintf("SubFeature %d deleted", subFeatureID), "SubFeature deleted successfully", http.StatusOK))
}

// Policy handlers
func (a *ACL) addPolicyHandler(c echo.Context) error {
	policy := new(Policy)
	if err := c.Bind(policy); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.AddPolicy(c.Request().Context(), policy); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusCreated, goresponse.GenerateSingleDataResponse(policy, "Policy added successfully", http.StatusCreated))
}

func (a *ACL) removePolicyHandler(c echo.Context) error {
	policyID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.RemovePolicy(c.Request().Context(), policyID); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(fmt.Sprintf("Policy %d removed", policyID), "Policy removed successfully", http.StatusOK))
}
