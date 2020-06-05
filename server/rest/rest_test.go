package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testHost = "1.1.1.1"
	testPort = 80
)

func TestPing(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := setupRouter(testHost, testPort)
	rec := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, fmt.Sprintf("PONGv%s %s %d", VERSION, testHost, testPort), rec.Body.String())


}