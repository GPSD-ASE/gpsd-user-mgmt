package tests

import (
	"gpsd-user-mgmt/src/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func TestPingRoute(t *testing.T) {
	router := gin.Default()

	router.GET("/api/v1/list", user.List)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/list", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"abc\":\"123\",\"qwe\":\"121\",\"zxc\":\"134\"}", w.Body.String())
}
