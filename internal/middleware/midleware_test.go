package middleware

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	tests := []struct {
		name                  string
		acceptEncoding        string
		expectContentEncoding string
		body                  string
	}{
		{"With Gzip", "gzip", "gzip", "Test 1"},
		{"Without Gzip", "", "", "Test 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.body))
				w.WriteHeader(http.StatusOK)
			})
			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			zb.Write([]byte(tt.body))
			zb.Close()

			r := httptest.NewRequest(http.MethodPost, "/test", buf)
			if tt.acceptEncoding != "" {
				r.Header.Set("Accept-Encoding", tt.acceptEncoding)
				r.Header.Set("Content-Encoding", tt.acceptEncoding)
			}

			rw := httptest.NewRecorder()
			m := GzipMiddleware(testHandler)
			m.ServeHTTP(rw, r)
			res := rw.Result()
			defer res.Body.Close()
			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

func TestHashMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test 1"))
		w.WriteHeader(http.StatusOK)
	})
	h := hmac.New(sha256.New, []byte("key"))
	h.Write([]byte("Test 1"))
	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("Test 1"))
	r.Header.Add("HashSHA256", hex.EncodeToString(h.Sum(nil)))

	rw := httptest.NewRecorder()
	HashMiddleware("key")(testHandler).ServeHTTP(rw, r)
	res := rw.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
