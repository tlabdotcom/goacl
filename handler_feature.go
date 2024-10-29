package goacl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tlabdotcom/gohelper"
	"github.com/tlabdotcom/goresponse"
)

// Feature handlers
func (a *ACL) listFeatureHandler(c echo.Context) error {
	datas, err := a.ListFeatures(c.Request().Context())
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

// func (a *ACL) updateFeatureHandler(c echo.Context) error {
// 	feature := new(Feature)
// 	if err := c.Bind(feature); err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
// 	}
// 	var err error
// 	feature.ID, err = gohelper.StringToInt64(c.Param("id"))
// 	if err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
// 	}
// 	if err := a.UpdateFeature(c.Request().Context(), feature); err != nil {
// 		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
// 	}
// 	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(feature, "Feature updated successfully", http.StatusOK))
// }

func (a *ACL) detailFeatureHandler(c echo.Context) error {
	featureID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	data, err := a.DetailFeature(c.Request().Context(), featureID)
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(data, "Get feature by id successfully", http.StatusOK))
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
