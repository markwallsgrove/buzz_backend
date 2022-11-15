package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/markwallsgrove/muzz_devops/src/routes"
	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMETextHTML)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	index := routes.IndexController{}
	assert.NoError(t, index.Index(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!!", rec.Body.String())
}
