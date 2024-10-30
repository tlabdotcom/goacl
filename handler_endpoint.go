package goacl

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tlabdotcom/gohelper"
	"github.com/tlabdotcom/goresponse"
)

func (a *ACL) deleteEndpointHandler(c echo.Context) error {
	endpointID, err := gohelper.StringToInt64(c.Param("id"))
	if err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusUnprocessableEntity).AddError(err).JSON(c)
	}
	if err := a.DeleteEndpoint(c.Request().Context(), endpointID); err != nil {
		return goresponse.NewStandardErrorResponse(http.StatusInternalServerError).AddError(err).JSON(c)
	}
	return c.JSON(http.StatusOK, goresponse.GenerateSingleDataResponse(fmt.Sprintf("Endpoint %d deleted", endpointID), "Endpoint deleted successfully", http.StatusOK))
}
