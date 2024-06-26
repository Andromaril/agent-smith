// Package gzip поддержку сжатия и декомпрессии данных
package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andromaril/agent-smith/internal/errormetric"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter создание нового compressWriter
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
		form, err := c.zw.Write(p)
		e := errormetric.NewMetricError(err)
		return form, fmt.Errorf("error gzip %w", e)
	} else {
		c.zw = nil
		return c.w.Write(p)
	}
}

func (c *compressWriter) WriteHeader(statusCode int) {
	contentType := c.w.Header().Get("Content-Type")
	support := strings.Contains(contentType, "application/json")
	support2 := strings.Contains(contentType, "text/html")
	if support || support2 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if c.zw != nil {
		return c.zw.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader создание нового ридера
func NewCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return nil, fmt.Errorf("error %w", e)
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
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error %w", e)
	}
	return c.zr.Close()
}
