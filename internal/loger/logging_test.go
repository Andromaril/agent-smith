package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithLogging(t *testing.T) {
	sugar := zap.NewExample().Sugar()
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test 1"))
		w.WriteHeader(http.StatusOK)
	})
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	rw := httptest.NewRecorder()
	WithLogging(*sugar)(testHandler).ServeHTTP(rw, r)
	res := rw.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
