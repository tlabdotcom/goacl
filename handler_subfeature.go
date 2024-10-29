package goacl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tlabdotcom/gohelper"
	"github.com/tlabdotcom/goresponse"
)

// SubFeature handlers
// func (a *ACL) listSubFeatureHandler(c echo.Context) error {
// 	datas, err := a.ListSubFeatures(c.Request().Context())
// 	if err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
// 	}
// 	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(datas, "success", http.StatusOK))
// }

// func (a *ACL) createSubFeatureHandler(c echo.Context) error {
// 	subFeature := new(SubFeature)
// 	if err := c.Bind(subFeature); err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
// 	}
// 	if err := a.CreateSubFeature(c.Request().Context(), subFeature); err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
// 	}
// 	return c.JSON(http.StatusCreated, goresponse.GenerateSingleDataResponse(subFeature, "SubFeature created successfully", http.StatusCreated))
// }

func (a *ACL) updateSubFeatureHandler(c echo.Context) error {
	params := new(SubFeatureParam)
	if err := c.Bind(params); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.UpdateSubFeature(c.Request().Context(), params); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}

	data, err := a.GetSubFeatureByID(c.Request().Context(), params.ID)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(data, "SubFeature updated successfully", http.StatusOK))
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
