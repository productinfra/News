package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": "just a test"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check whether the response contains the expected login error
	// 1. Method 1: Check if the response contains the specified string
	//assert.Equal(t, w.Body.String(), "Login required")

	// 2. Method 2: Deserialize the response content into ResponseData and check if the fields match expectations
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err: %v\n", err)
	}
	assert.Equal(t, res.Code, CodeNotLogin)
}
