package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/apierror"
	"github.com/syukur91/ischool-monitor/pkg/middleware"
	"github.com/syukur91/ischool-monitor/pkg/query"
	"github.com/syukur91/ischool-monitor/pkg/response"
	"github.com/syukur91/ischool-monitor/service"
)

// Wali_KelasHandler ...
type Wali_KelasHandler struct {
	Wali_KelasService *service.Wali_KelasService
}

// SetRoutes ...
func (h *Wali_KelasHandler) SetRoutes(r *echo.Group) {
	r.POST("/wali_kelass", h.createWali_Kelas)
	r.POST("/wali_kelass-grid", h.gridWali_Kelass, middleware.KendoGrid)
	r.GET("/wali_kelass/:id", h.getWali_Kelas)
	r.POST("/wali_kelass/:id", h.updateWali_Kelas)
	r.DELETE("/wali_kelass/:id", h.deleteWali_Kelas)
}

func (h *Wali_KelasHandler) createWali_Kelas(c echo.Context) error {
	createWali_Kelas := new(schema.CreateWali_KelasRequest)
	err := c.Bind(createWali_Kelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get wali_kelas data. Probably content-type is not match with actual body type", errors.New("createWali_Kelas: Failed to get wali_kelas data"))
	}

	err = c.Validate(createWali_Kelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Wali_Kelas data invalid. One or more required fields is not set", errors.New("createWali_Kelas: invalid wali_kelas data"))
	}

	createWali_KelasResponse, err := h.Wali_KelasService.CreateWali_Kelas(createWali_Kelas)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, createWali_KelasResponse)
}

func (h *Wali_KelasHandler) gridWali_Kelass(c echo.Context) error {
	gridParams := c.Get("gridParams").(*query.GridParams)

	data, count, err := h.Wali_KelasService.ListWali_Kelass(gridParams)
	if err != nil {
		return err
	}

	return response.JSONGrid(c, http.StatusOK, data, len(data), count)
}

func (h *Wali_KelasHandler) getWali_Kelas(c echo.Context) error {
	id := c.Param("id")

	getWali_KelasResponse, err := h.Wali_KelasService.GetWali_Kelas(id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, getWali_KelasResponse)
}

func (h *Wali_KelasHandler) updateWali_Kelas(c echo.Context) error {
	id := c.Param("id")

	updateWali_Kelas := new(schema.UpdateWali_KelasRequest)
	err := c.Bind(updateWali_Kelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get wali_kelas data. Probably content-type is not match with actual body type", errors.New("createWali_Kelas: Failed to get wali_kelas data"))
	}

	updateWali_KelasResponse, err := h.Wali_KelasService.UpdateWali_Kelas(id, updateWali_Kelas)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, updateWali_KelasResponse)
}
func (h *Wali_KelasHandler) deleteWali_Kelas(c echo.Context) error {
	id := c.Param("id")

	err := h.Wali_KelasService.DeleteWali_Kelas(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
