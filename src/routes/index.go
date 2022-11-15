package routes

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// IndexController friendly index page
type IndexController struct {
}

// Index handle the index route
func (i *IndexController) Index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!!")
}
