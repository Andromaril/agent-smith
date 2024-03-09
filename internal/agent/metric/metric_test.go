package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSendGaugeMetric(t *testing.T) {
	s := storage.NewMemStorage(serverflag.StoreInterval == 0, serverflag.FileStoragePath)
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{pattern}/{name}/{value}", handler.GaugeandCounter(s))

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid gauge", url: "/update/gauge/test/1", expectedCode: http.StatusOK},
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

func TestSendCounterMetric(t *testing.T) {
	s := storage.NewMemStorage(serverflag.StoreInterval == 0, serverflag.FileStoragePath)
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{pattern}/{name}/{value}", handler.GaugeandCounter(s))

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid counter", url: "/update/counter/test/2", expectedCode: http.StatusOK},
		{name: "unvalid url", url: "/update/counter/", expectedCode: http.StatusNotFound},
		{name: "unvalid method", url: "/update/counter/test/rt", expectedCode: http.StatusBadRequest},
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
