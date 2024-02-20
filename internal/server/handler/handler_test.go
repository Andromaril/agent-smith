package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andromaril/agent-smith/internal/server/storage"
    "github.com/stretchr/testify/assert"
)

func TestGaugeandCounter(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedCode   int
	}{
		{name: "valid gauge", url: "/update/gauge/test/1", expectedCode: http.StatusOK},
        {name: "valid counter", url: "/update/gauge/test/2", expectedCode: http.StatusOK},
        {name: "unvalid url", url: "/update/gauge/", expectedCode: http.StatusNotFound},
        {name: "unvalid method", url: "/update/gauge/test/rt", expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			s := storage.NewMemStorage()
			h := http.HandlerFunc(GaugeandCounter(s))
			h(w, r)
            assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
