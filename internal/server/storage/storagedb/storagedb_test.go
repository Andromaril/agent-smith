// Package storagedb необходим для работы с базой данных, где хранятся метрики

package storagedb

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andromaril/agent-smith/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestStorageDB_Bootstrap(t *testing.T) {
	//var ctx context.Background()
	//ctx := context.Background()
	db, mock, err := sqlmock.New()
	//s := StorageDB{DB: db}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`
	CREATE TABLE IF NOT EXISTS gauge (
		id SERIAL PRIMARY KEY,
		key varchar(100) UNIQUE NOT NULL,
		value DOUBLE PRECISION NOT NULL
		);
	`)).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	CREATE TABLE IF NOT EXISTS counter (
		id SERIAL PRIMARY KEY,
		key varchar(100) UNIQUE NOT NULL,
		value bigint
		);
	`)).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
	// _, err = s.Init("postgres://postgres:qwerty123@localhost:5432/gr", ctx)
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestStorageDB_CounterAndGaugeUpdateMetrics(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
	`)).WithArgs(
		"gauge1",
		1.1,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO counter (key, value)
			VALUES($1, $2)
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
	`)).WithArgs(
		"counter1",
		1,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
	`)).WithArgs(
		"gauge2",
		2.2,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO counter (key, value)
			VALUES($1, $2)
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
	`)).WithArgs(
		"counter2",
		1,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
	list := make([]model.Gauge, 5)
	list2 := make([]model.Counter, 5)
	list = append(list, model.Gauge{Key: "gauge1", Value: 1.1})
	list = append(list, model.Gauge{Key: "gauge2", Value: 2.2})
	list2 = append(list2, model.Counter{Key: "counter1", Value: 1})
	list2 = append(list2, model.Counter{Key: "counter2", Value: 2})
	err = s.CounterAndGaugeUpdateMetrics(list, list2)
	if err != nil {
		t.Error(err)
	}

}

func TestStorageDB_NewGauge(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
	`)).WithArgs(
		"gauge1",
		1.1,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
	`)).WithArgs(
		"gauge2",
		2.2,
	).WillReturnResult(driver.ResultNoRows)
	err1 := s.NewGauge("gauge1", 1.1)
	if err1 != nil {
		t.Error(err1)
	}
	err2 := s.NewGauge("gauge2", 2.2)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestStorageDB_NewCounter(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO counter (key, value)
	VALUES($1, $2) 
	ON CONFLICT (key) 
	DO UPDATE SET value = counter.value + $2;
	`)).WithArgs(
		"counter1",
		1,
	).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
	INSERT INTO counter (key, value)
	VALUES($1, $2) 
	ON CONFLICT (key) 
	DO UPDATE SET value = counter.value + $2;
	`)).WithArgs(
		"counter2",
		2,
	).WillReturnResult(driver.ResultNoRows)
	err1 := s.NewCounter("counter1", 1)
	if err1 != nil {
		t.Error(err1)
	}
	err2 := s.NewCounter("counter2", 2)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestStorageDB_GetCounter(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM counter WHERE key=$1`)).
		WithArgs("counter1").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))
	k, err := s.GetCounter("counter1")
	if k != 1 {
		t.Errorf("unexpected value, expected %v", k)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestStorageDB_GetGauge(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM gauge WHERE key=$1`)).
		WithArgs("gauge1").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1.1))
	k, err := s.GetGauge("gauge1")
	if k != 1.1 {
		t.Errorf("unexpected value, expected %v", k)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestStorageDB_GetIntMetric(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM counter`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"counter1", 1}).
			AddRows([]driver.Value{"counter2", 2}))
	k, err := s.GetIntMetric()
	if k["counter1"] != 1 && k["counter2"] != 2 {
		t.Errorf("unexpected value, expected %v", k)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestStorageDB_GetFloatMetric(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM gauge`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"gauge1", 1.1}).
			AddRows([]driver.Value{"gauge2", 2.2}))
	k, err := s.GetFloatMetric()
	if k["gauge1"] != 1.1 && k["gauge2"] != 2.1 {
		t.Errorf("unexpected value, expected %v", k)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestStorageDB_Ping(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = s.Ping()
	assert.NoError(t, err)
}

func TestStorageDB_Load(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx, Path: "test"}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = s.Load(s.Path)
	assert.NoError(t, err)
}

func TestStorageDB_Save(t *testing.T) {
	ctx := context.Background()
	db, _, err := sqlmock.New()
	s := StorageDB{DB: db, Ctx: ctx, Path: "test"}
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	err = s.Save(s.Path)
	assert.NoError(t, err)
}
