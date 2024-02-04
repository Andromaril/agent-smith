package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGaugeandCounter(t *testing.T) {
	s := storage.NewMemStorage()
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{change}/{name}/{value}", GaugeandCounter(s))

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid gauge", url: "/update/gauge/test/1", expectedCode: http.StatusOK},
		{name: "valid counter", url: "/update/gauge/test/2", expectedCode: http.StatusOK},
		{name: "unvalid url", url: "/update/gauge/", expectedCode: http.StatusNotFound},
		{name: "unvalid method", url: "/update/gauge/test/rt", expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, err := http.NewRequest(http.MethodPost, r.URL+tt.url, nil)
			assert.NoError(t, err)
			response, err := r.Client().Do(r1)
			assert.NoError(t, err)
			response.Body.Close()
			assert.Equal(t, tt.expectedCode, response.StatusCode)
		})
	}
}

func TestGetMetric(t *testing.T) {
	s := storage.NewMemStorage()
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{change}/{name}/{value}", GaugeandCounter(s))
	ts.Get("/value/{change}/{name}", GetMetric(s))

	tests := []struct {
		name         string
		url1         string
		url2         string
		expectedCode int
	}{
		{name: "valid gauge", url1: "/update/gauge/test/1", url2: "/value/gauge/test", expectedCode: http.StatusOK},
		{name: "valid counter", url1: "/update/gauge/test/2", url2: "/value/gauge/test", expectedCode: http.StatusOK},
		{name: "unvalid url", url1: "/update/gauge/", url2: "/value/gauge", expectedCode: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err1 := http.NewRequest(http.MethodPost, r.URL+tt.url1, nil)
			r1, err2 := http.NewRequest(http.MethodGet, r.URL+tt.url2, nil)
			assert.NoError(t, err1)
			assert.NoError(t, err2)
			response, err := r.Client().Do(r1)
			assert.NoError(t, err)
			response.Body.Close()
			assert.Equal(t, tt.expectedCode, response.StatusCode)
		})
	}
}

func TestGetHtmlMetric(t *testing.T) {
	s := storage.NewMemStorage()
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Get("/", GetHtmlMetric(s))
	r1, err1 := http.NewRequest(http.MethodGet, r.URL, nil)
	assert.NoError(t, err1)
	response, err2 := r.Client().Do(r1)
	assert.NoError(t, err2)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/html", response.Header.Get("Content-Type"))
}
