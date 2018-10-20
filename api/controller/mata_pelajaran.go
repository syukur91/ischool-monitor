package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"gitlab.com/nextid/tenant-api/api/schema"
	"gitlab.com/nextid/tenant-api/pkg/apierror"
	"gitlab.com/nextid/tenant-api/pkg/middleware"
	"gitlab.com/nextid/tenant-api/pkg/query"
	"gitlab.com/nextid/tenant-api/pkg/response"
	"gitlab.com/nextid/tenant-api/service"
)

// Mata_PelajaranHandler ...
type Mata_PelajaranHandler struct {
	Mata_PelajaranService *service.Mata_PelajaranService
}

// SetRoutes ...
func (h *Mata_PelajaranHandler) SetRoutes(r *echo.Group) {
	r.POST("/mata_pelajarans", h.createMata_Pelajaran)
	r.POST("/mata_pelajarans-grid", h.gridMata_Pelajarans, middleware.KendoGrid)
	r.GET("/mata_pelajarans/:id", h.getMata_Pelajaran)
	r.POST("/mata_pelajarans/:id", h.updateMata_Pelajaran)
	r.DELETE("/mata_pelajarans/:id", h.deleteMata_Pelajaran)
}

func (h *Mata_PelajaranHandler) createMata_Pelajaran(c echo.Context) error {
	createMata_Pelajaran := new(schema.CreateMata_PelajaranRequest)
	err := c.Bind(createMata_Pelajaran)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get mata_pelajaran data. Probably content-type is not match with actual body type", errors.New("createMata_Pelajaran: Failed to get mata_pelajaran data"))
	}

	err = c.Validate(createMata_Pelajaran)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Mata_Pelajaran data invalid. One or more required fields is not set", errors.New("createMata_Pelajaran: invalid mata_pelajaran data"))
	}

	createMata_PelajaranResponse, err := h.Mata_PelajaranService.CreateMata_Pelajaran(createMata_Pelajaran)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, createMata_PelajaranResponse)
}

func (h *Mata_PelajaranHandler) gridMata_Pelajarans(c echo.Context) error {
	gridParams := c.Get("gridParams").(*query.GridParams)

	data, count, err := h.Mata_PelajaranService.ListMata_Pelajarans(gridParams)
	if err != nil {
		return err
	}

	return response.JSONGrid(c, http.StatusOK, data, len(data), count)
}

func (h *Mata_PelajaranHandler) getMata_Pelajaran(c echo.Context) error {
	id := c.Param("id")

	getMata_PelajaranResponse, err := h.Mata_PelajaranService.GetMata_Pelajaran(id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, getMata_PelajaranResponse)
}

func (h *Mata_PelajaranHandler) updateMata_Pelajaran(c echo.Context) error {
	id := c.Param("id")

	updateMata_Pelajaran := new(schema.UpdateMata_PelajaranRequest)
	err := c.Bind(updateMata_Pelajaran)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get mata_pelajaran data. Probably content-type is not match with actual body type", errors.New("createMata_Pelajaran: Failed to get mata_pelajaran data"))
	}

	updateMata_PelajaranResponse, err := h.Mata_PelajaranService.UpdateMata_Pelajaran(id, updateMata_Pelajaran)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, updateMata_PelajaranResponse)
}
func (h *Mata_PelajaranHandler) deleteMata_Pelajaran(c echo.Context) error {
	id := c.Param("id")

	err := h.Mata_PelajaranService.DeleteMata_Pelajaran(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
