package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetIndex(t *testing.T) {
	router := gin.New()
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	router.GET("/", GetIndex)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestPatchConfig(t *testing.T) {
	router := gin.New()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("PATCH", "/", strings.NewReader("{}"))

	router.PATCH("/", PatchConfig)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
