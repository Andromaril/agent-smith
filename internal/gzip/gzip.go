package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	contentType := c.w.Header().Get("Content-Type")
	support := strings.Contains(contentType, "application/json")
	support2 := strings.Contains(contentType, "text/html")
	if support || support2 {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.zw = gzip.NewWriter(c.w)
		buf, err := c.zw.Write(p)
		// err2 := c.Close()
		// if err2 != nil {
		// 	return 0, err2
		// }

		return buf, err
	} else {
		c.zw = nil
		return c.w.Write(p)
	}
}

func (c *compressWriter) WriteHeader(statusCode int) {
	// contentType := c.w.Header().Get("Content-Type")
	// support := strings.Contains(contentType, "application/json")
	// support2 := strings.Contains(contentType, "text/html")
	// if support || support2 {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	//if c.zw != nil {
	return c.zw.Close()
	//}
	//return nil
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
