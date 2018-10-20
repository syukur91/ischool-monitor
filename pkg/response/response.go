package response

import (
	"github.com/labstack/echo"
)

// Response JSONAPI object
type Response struct {
	Errors []Error     `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Links  interface{} `json:"links,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Total  int         `json:"total,omitempty"`
}

// Error object
type Error struct {
	Status int    `json:"status,omitempty"`
	Code   int    `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// JSON is
func JSON(c echo.Context, status int, data interface{}) error {
	r := new(Response)
	r.Data = data
	return c.JSON(status, r)
}

// JSONGrid is
func JSONGrid(c echo.Context, status int, data interface{}, length int, count int) error {
	r := struct {
		Data  interface{} `json:"data"`
		Total int         `json:"total"`
	}{
		make([]interface{}, 0),
		count,
	}

	if length > 0 {
		r.Data = data
	}

	return c.JSON(status, r)
}
