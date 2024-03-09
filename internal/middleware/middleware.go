package middleware

import (
	"net/http"
	"strings"

	"github.com/andromaril/agent-smith/internal/gzip"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := r.Header.Get("Content-Type")
		support := strings.Contains(contentType, "application/json")
		support2 := strings.Contains(contentType, "text/html")
		if support || support2 {
			if supportsGzip {
				cw := gzip.NewCompressWriter(w)
				ow = cw
				ow.Header().Set("Content-Encoding", "gzip")
				defer cw.Close()
			}
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	})

}
