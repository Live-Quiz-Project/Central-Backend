package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserIntegration(t *testing.T) {
	// Set GIN to Test Mode
	gin.SetMode(gin.TestMode)

	// Prepare the GIN router
	router := gin.Default()
	router.POST("/v1/user", )

	// Prepare the request body
	requestBody := map[string]interface{}{
		"email":    "test1@somemail.com",
		"password": "test123",
		"name":     "TestProfile 1",
	}
	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Create a request with the prepared body
	req, err := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request to the GIN router
	router.ServeHTTP(rr, req)

	// Verify the status code and response body
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201, got %v", rr.Code)

	// You can also check the response body or headers if needed
	// responseBody := rr.Body.String()
	// assert.Contains(t, responseBody, "some expected content", "Expected response body to contain 'some expected content'")
}
