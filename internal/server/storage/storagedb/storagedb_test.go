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
)

func TestStorageDB_Bootstrap(t *testing.T) {
	//var ctx context.Background()
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	s := StorageDB{DB: db}
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
	//s.DB, err = sql.Open("pgx", "test")
	_, err = s.Init("postgres://postgres:qwerty123@localhost:5432/gr", ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestStorageDB_CounterAndGaugeUpdateMetrics(t *testing.T) {
	
}
