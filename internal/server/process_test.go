package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	storeschema "github.com/compliance-framework/configuration-service/internal/stores/schema"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterProcess(t *testing.T) {
	f := FakeDriver{}
	s := &Server{Driver: f}
	p := echo.New()
	err := s.RegisterProcess(p)
	assert.Nil(t, err)
	expected := map[string]bool{
		"GET/job-results/:uuid": false,
		"GET/job-results":       false,
	}
	for _, routes := range p.Routes() {
		t := fmt.Sprintf("%s%s", routes.Method, routes.Path)
		if _, ok := expected[t]; ok {
			expected[t] = true
		}
	}
	for k, v := range expected {
		assert.True(t, v, fmt.Sprintf("expected route %s not found", k))
	}
}

func TestGetJobResult(t *testing.T) {
	testCases := []struct {
		name         string
		getFn        func(id string, object interface{}) error
		path         string
		params       map[string]string
		requestPath  string
		expectedCode int
	}{
		{
			name: "get-job-result",
			getFn: func(id string, object interface{}) error {
				// Simulate a successful Get call here
				return nil
			},
			path:         "/job-results/:uuid",
			params:       map[string]string{"uuid": "1234"},
			requestPath:  "/job-results/1234",
			expectedCode: 200,
		},
		{
			name: "get-job-result-not-found",
			getFn: func(id string, object interface{}) error {
				return storeschema.NotFoundErr{}
			},
			path:         "/job-results/:uuid",
			params:       map[string]string{"uuid": "1236"},
			requestPath:  "/job-results/1236",
			expectedCode: 404,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := FakeDriver{}
			f.GetFn = tc.getFn

			s := &Server{Driver: f}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := s.GetJobResult(c)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedCode, c.Response().Status)
		})
	}
}

func TestGetJobResults(t *testing.T) {
	testCases := []struct {
		name     string
		getAllFn func(id string, object interface{}) ([]interface{}, error)
		path     string

		requestPath  string
		expectedCode int
	}{
		{
			name: "get-job-results",
			getAllFn: func(id string, object interface{}) ([]interface{}, error) {
				// Simulate a successful Get call here
				return nil, nil
			},
			path: "/job-results",

			requestPath:  "/job-results",
			expectedCode: 200,
		},
		{
			name: "get-job-result-not-found",
			getAllFn: func(id string, object interface{}) ([]interface{}, error) {
				// Simulate a successful Get call here
				return nil, fmt.Errorf("boom")
			},
			path: "/job-results",

			requestPath:  "/job-results",
			expectedCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := FakeDriver{}
			f.GetAllFn = tc.getAllFn

			s := &Server{Driver: f}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := s.GetJobResults(c)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedCode, c.Response().Status)
		})
	}
}
