package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestSendGaugeMetric(t *testing.T) {

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid gauge", url: "http://localhost:8080/update/gauge/test/1", expectedCode: http.StatusOK},
		{name: "unvalid url", url: "http://localhost:8080/update/gauge/", expectedCode: http.StatusNotFound},
		{name: "unvalid method", url: "http://localhost:8080/update/gauge/test/rt", expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			s := storage.NewMemStorage()
			h := http.HandlerFunc(handler.GaugeandCounter(s))
			h(w, r)
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestSendCounterMetric(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid counter", url: "http://localhost:8080/update/counter/test/2", expectedCode: http.StatusOK},
		{name: "unvalid url", url: "http://localhost:8080/update/counter/", expectedCode: http.StatusNotFound},
		{name: "unvalid method", url: "http://localhost:8080/update/counter/test/rt", expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			s := storage.NewMemStorage()
			h := http.HandlerFunc(handler.GaugeandCounter(s))
			h(w, r)
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
