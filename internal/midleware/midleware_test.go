package midleware

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/andromaril/agent-smith/internal/errormetric"
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

func TestCryptoMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test 1"))
		w.WriteHeader(http.StatusOK)
	})
	data, err := os.ReadFile("../../key/key.pem.pub")
	if err != nil {
		e := errormetric.NewMetricError(err)
		t.Error("error read file %w", e)
	}
	pemDecode, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(pemDecode.Bytes)
	if err != nil {
		e := errormetric.NewMetricError(err)
		t.Error("error decode public key %w", e)
	}
	var buf2 []byte
	buf2, err = rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte("Test 1"))
	if err != nil {
		e := errormetric.NewMetricError(err)
		t.Error("error decode public key %w", e)
	}
	data2, err := os.ReadFile("../../key/sever.key")
	if err != nil {
		e := errormetric.NewMetricError(err)
		t.Error("error read file %w", e)
	}
	pemDecode2, _ := pem.Decode(data2)
	priv, err := x509.ParsePKCS1PrivateKey(pemDecode2.Bytes)

	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(string(buf2)))

	rw := httptest.NewRecorder()
	CryptoMiddleware(priv)(testHandler).ServeHTTP(rw, r)
	res := rw.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
