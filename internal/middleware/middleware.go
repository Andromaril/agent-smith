package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/andromaril/agent-smith/internal/gzip"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := ow.Header().Values("Content-Type")
		//accept := r.Header.Get("Accept")
		support := slices.Contains(contentType, "application/json")
		support2 := slices.Contains(contentType, "text/html")
		fmt.Print(support, support2)
		//support3 := strings.Contains(accept, "html/text")
		//if support || support2 {
		//По заданию мы проверяем не только Accept-Encoding, но и Content-Type для принятия решения о сжатии ответа.
		if supportsGzip {
			//if support || support2 {
			cw := gzip.NewCompressWriter(w)
			ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
			//}
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		//if support || support2 {
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
		//}
		h.ServeHTTP(ow, r)
	})

}
