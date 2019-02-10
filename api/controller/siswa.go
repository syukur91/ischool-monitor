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

// SiswaHandler ...
type SiswaHandler struct {
	SiswaService *service.SiswaService
}

// SetRoutes ...
func (h *SiswaHandler) SetRoutes(r *echo.Group) {
	r.POST("/siswas", h.createSiswa)
	r.POST("/siswas-grid", h.gridSiswas, middleware.KendoGrid)
	r.GET("/siswas/:id", h.getSiswa)
	r.POST("/siswas/:id", h.updateSiswa)
	r.DELETE("/siswas/:id", h.deleteSiswa)
}

func (h *SiswaHandler) createSiswa(c echo.Context) error {
	createSiswa := new(schema.CreateSiswaRequest)
	err := c.Bind(createSiswa)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get siswa data. Probably content-type is not match with actual body type", errors.New("createSiswa: Failed to get siswa data"))
	}

	err = c.Validate(createSiswa)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Siswa data invalid. One or more required fields is not set", errors.New("createSiswa: invalid siswa data"))
	}

	createSiswaResponse, err := h.SiswaService.CreateSiswa(createSiswa)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, createSiswaResponse)
}

func (h *SiswaHandler) gridSiswas(c echo.Context) error {
	gridParams := c.Get("gridParams").(*query.GridParams)

	data, count, err := h.SiswaService.ListSiswas(gridParams)
	if err != nil {
		return err
	}

	return response.JSONGrid(c, http.StatusOK, data, len(data), count)
}

func (h *SiswaHandler) getSiswa(c echo.Context) error {
	id := c.Param("id")

	getSiswaResponse, err := h.SiswaService.GetSiswa(id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, getSiswaResponse)
}

func (h *SiswaHandler) updateSiswa(c echo.Context) error {
	id := c.Param("id")

	updateSiswa := new(schema.UpdateSiswaRequest)
	err := c.Bind(updateSiswa)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get siswa data. Probably content-type is not match with actual body type", errors.New("createSiswa: Failed to get siswa data"))
	}

	updateSiswaResponse, err := h.SiswaService.UpdateSiswa(id, updateSiswa)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, updateSiswaResponse)
}
func (h *SiswaHandler) deleteSiswa(c echo.Context) error {
	id := c.Param("id")

	err := h.SiswaService.DeleteSiswa(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
