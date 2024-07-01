// Package midleware содержит middleware для сжатия и хеширования
package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/andromaril/agent-smith/internal/gzip"
)

// GzipMiddleware middleware для передачи сжатия данных запроса
func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := gzip.NewCompressWriter(w)
			ow = cw
			//ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//ow.Header().Set("Content-Encoding", "gzip")
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	})

}

// HashMiddleware middleware для хеширования
func HashMiddleware(key string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("HashSHA256") != "" {
				hash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				body, err1 := io.ReadAll(r.Body)
				if err1 != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				h := hmac.New(sha256.New, []byte(key))
				h.Write(body)
				if !hmac.Equal(h.Sum(nil), hash) {
					w.WriteHeader(http.StatusBadRequest)
				}
			}
			h.ServeHTTP(w, r)

		})

	}
}

func CryptoMiddleware(key *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			body, err1 := io.ReadAll(r.Body)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			msgDecode, err := rsa.DecryptPKCS1v15(rand.Reader, key, body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(msgDecode))

			h.ServeHTTP(w, r)

		})

	}
}

func IPMiddleware(cidrRange string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			_, ipnet, err := net.ParseCIDR(cidrRange)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ipStr := r.Header.Get("X-Real-IP")
			// парсим ip
			ip := net.ParseIP(ipStr)
			ipcheck := ipnet.Contains(ip)
			if !ipcheck {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)

		})

	}
}
