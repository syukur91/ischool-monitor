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

// KelasHandler ...
type KelasHandler struct {
	KelasService *service.KelasService
}

// SetRoutes ...
func (h *KelasHandler) SetRoutes(r *echo.Group) {
	r.POST("/kelass", h.createKelas)
	r.POST("/kelass-grid", h.gridKelass, middleware.KendoGrid)
	r.GET("/kelass/:id", h.getKelas)
	r.POST("/kelass/:id", h.updateKelas)
	r.DELETE("/kelass/:id", h.deleteKelas)
}

func (h *KelasHandler) createKelas(c echo.Context) error {
	createKelas := new(schema.CreateKelasRequest)
	err := c.Bind(createKelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get kelas data. Probably content-type is not match with actual body type", errors.New("createKelas: Failed to get kelas data"))
	}

	err = c.Validate(createKelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Kelas data invalid. One or more required fields is not set", errors.New("createKelas: invalid kelas data"))
	}

	createKelasResponse, err := h.KelasService.CreateKelas(createKelas)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, createKelasResponse)
}

func (h *KelasHandler) gridKelass(c echo.Context) error {
	gridParams := c.Get("gridParams").(*query.GridParams)

	data, count, err := h.KelasService.ListKelass(gridParams)
	if err != nil {
		return err
	}

	return response.JSONGrid(c, http.StatusOK, data, len(data), count)
}

func (h *KelasHandler) getKelas(c echo.Context) error {
	id := c.Param("id")

	getKelasResponse, err := h.KelasService.GetKelas(id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, getKelasResponse)
}

func (h *KelasHandler) updateKelas(c echo.Context) error {
	id := c.Param("id")

	updateKelas := new(schema.UpdateKelasRequest)
	err := c.Bind(updateKelas)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get kelas data. Probably content-type is not match with actual body type", errors.New("createKelas: Failed to get kelas data"))
	}

	updateKelasResponse, err := h.KelasService.UpdateKelas(id, updateKelas)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, updateKelasResponse)
}
func (h *KelasHandler) deleteKelas(c echo.Context) error {
	id := c.Param("id")

	err := h.KelasService.DeleteKelas(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
