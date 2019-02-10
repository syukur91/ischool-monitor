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

// UserHandler ...
type UserHandler struct {
	UserService *service.UserService
}

// SetRoutes ...
func (h *UserHandler) SetRoutes(r *echo.Group) {
	r.POST("/users", h.createUser)
	r.POST("/users-grid", h.gridUsers, middleware.KendoGrid)
	r.GET("/users/:id", h.getUser)
	r.POST("/users/:id", h.updateUser)
	r.DELETE("/users/:id", h.deleteUser)
}

func (h *UserHandler) createUser(c echo.Context) error {
	createUser := new(schema.CreateUserRequest)
	err := c.Bind(createUser)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get user data. Probably content-type is not match with actual body type", errors.New("createUser: Failed to get user data"))
	}

	err = c.Validate(createUser)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "User data invalid. One or more required fields is not set", errors.New("createUser: invalid user data"))
	}

	createUserResponse, err := h.UserService.CreateUser(createUser)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, createUserResponse)
}

func (h *UserHandler) gridUsers(c echo.Context) error {
	gridParams := c.Get("gridParams").(*query.GridParams)

	data, count, err := h.UserService.ListUsers(gridParams)
	if err != nil {
		return err
	}

	return response.JSONGrid(c, http.StatusOK, data, len(data), count)
}

func (h *UserHandler) getUser(c echo.Context) error {
	id := c.Param("id")

	getUserResponse, err := h.UserService.GetUser(id)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, getUserResponse)
}

func (h *UserHandler) updateUser(c echo.Context) error {
	id := c.Param("id")

	updateUser := new(schema.UpdateUserRequest)
	err := c.Bind(updateUser)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, http.StatusUnprocessableEntity, "Failed to get user data. Probably content-type is not match with actual body type", errors.New("createUser: Failed to get user data"))
	}

	updateUserResponse, err := h.UserService.UpdateUser(id, updateUser)
	if err != nil {
		return err
	}

	return response.JSON(c, http.StatusOK, updateUserResponse)
}
func (h *UserHandler) deleteUser(c echo.Context) error {
	id := c.Param("id")

	err := h.UserService.DeleteUser(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
