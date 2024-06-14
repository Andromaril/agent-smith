package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestGaugeandCounter(t *testing.T) {
	s := storage.NewMemStorage(false, "test")
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{pattern}/{name}/{value}", GaugeandCounter(s))

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{name: "valid gauge", url: "/update/gauge/test/1", expectedCode: http.StatusOK},
		{name: "valid counter", url: "/update/gauge/test/2", expectedCode: http.StatusOK},
		{name: "unvalid url", url: "/update/gauge/", expectedCode: http.StatusNotFound},
		{name: "unvalid method", url: "/update/gauge/test/rt", expectedCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, err := http.NewRequest(http.MethodPost, r.URL+tt.url, nil)
			assert.NoError(t, err)
			response, err := r.Client().Do(r1)
			assert.NoError(t, err)
			response.Body.Close()
			assert.Equal(t, tt.expectedCode, response.StatusCode)
		})
	}
}

func ExampleGaugeandCounter() {
	// Выполняет Post-запрос по адресу /update/ с url параметрами:
	// /update/{pattern}/{name}/{value}, где
	// pattern - тип метрики, gauge float64 или counter int64
	// name - имя метрики
	// value - значение
}

func TestGetMetric(t *testing.T) {
	s := storage.NewMemStorage(false, "test")
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/update/{pattern}/{name}/{value}", GaugeandCounter(s))
	ts.Get("/value/{pattern}/{name}", GetMetric(s))

	tests := []struct {
		name         string
		url1         string
		url2         string
		expectedCode int
	}{
		{name: "valid gauge", url1: "/update/gauge/test/1", url2: "/value/gauge/test", expectedCode: http.StatusOK},
		{name: "valid counter", url1: "/update/gauge/test/2", url2: "/value/gauge/test", expectedCode: http.StatusOK},
		{name: "unvalid url", url1: "/update/gauge/", url2: "/value/gauge", expectedCode: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, err1 := http.NewRequest(http.MethodPost, r.URL+tt.url1, nil)
			r2, err2 := http.NewRequest(http.MethodGet, r.URL+tt.url2, nil)
			assert.NoError(t, err1)
			assert.NoError(t, err2)
			response1, err4 := r.Client().Do(r1)
			response2, err5 := r.Client().Do(r2)
			assert.NoError(t, err4)
			assert.NoError(t, err5)
			response1.Body.Close()
			response2.Body.Close()
			assert.Equal(t, tt.expectedCode, response2.StatusCode)
		})
	}
}

func ExampleGetMetric() {
	// Выполняет Get-запрос по адресу /value/ с url параметрами:
	// /value/{pattern}/{name}/{value}, где
	// pattern - тип метрики, gauge float64 или counter int64
	// name - имя метрики
	// value - значение
}

func TestGetHTMLMetric(t *testing.T) {
	s := storage.NewMemStorage(false, "test")
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Get("/", GetHTMLMetric(s))
	r1, err1 := http.NewRequest(http.MethodGet, r.URL, nil)
	assert.NoError(t, err1)
	response, err2 := r.Client().Do(r1)
	assert.NoError(t, err2)
	response.Body.Close()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "text/html", response.Header.Get("Content-Type"))
}

func ExampleGetHTMLMetric() {
	// Выполняет Get-запрос по адресу /
	// В ответ выводит html-страницу со всеми метриками
}

func ExamplePing() {
	// Выполняет GET-запрос по адресу /ping и проверяет состояние базы данных
}

func TestPing(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	s := &storagedb.StorageDB{DB: db, Ctx: ctx}
	ts := chi.NewRouter()
	r := httptest.NewServer(ts)
	defer r.Close()
	ts.Post("/ping/", Ping(s))
	r1, err := http.NewRequest(http.MethodPost, r.URL+"/ping/", nil)
	assert.NoError(t, err)
	response, err := r.Client().Do(r1)
	assert.NoError(t, err)
	response.Body.Close()
}
