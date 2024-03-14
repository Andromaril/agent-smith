package middleware

import (
	"net/http"
	"strings"

	"github.com/andromaril/agent-smith/internal/gzip"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		//acceptEncoding := r.Header.Get("Accept-Encoding")
		//supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := ow.Header().Get("Content-Type")
		//accept := r.Header.Get("Accept")
		support := strings.Contains(contentType, "application/json")
		//support2 := slices.Contains(contentType, "text/html")
		//fmt.Print(support, support2)
		//support3 := strings.Contains(accept, "html/text")
		//if support {
		//По заданию мы проверяем не только Accept-Encoding, но и Content-Type для принятия решения о сжатии ответа.
		if support {
			//if support {
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
