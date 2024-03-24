package requester_test

import (
	"github.com/Waelson/go-o11y/internal/requester"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNormalize(t *testing.T) {
	httpReq := requester.NewHttpRequest()

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid URL",
			input:       "http://example.com?param=value",
			expected:    "http://example.com?param%3Dvalue",
			expectError: false,
		},
		{
			name:        "invalid URL",
			input:       "http://a b.com",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			normalized, err := httpReq.Normalize(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, normalized)
			}
		})
	}
}

func TestMakeRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpReq := requester.NewHttpRequest()

	mockURL := "http://example.com"
	mockResponse := `{"message": "hello world"}`
	httpmock.RegisterResponder("GET", mockURL, httpmock.NewStringResponder(http.StatusOK, mockResponse))

	// Valores esperados
	expectedBody := mockResponse
	expectedStatusCode := http.StatusOK

	body, statusCode, err := httpReq.MakeRequest(mockURL)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatusCode, statusCode)
	assert.Equal(t, expectedBody, body)

}
