package storagedb

import (
	"context"
	"database/sql"
)

type StorageDB struct {
	db   *sql.DB
	Path string
	ctx  context.Context
}

func (m *StorageDB) Init(path string, ctx context.Context) error {
	var err error
	m.ctx = ctx
	m.Path = path
	m.db, err = sql.Open("pgx", path)
	if err != nil {
		return err
	}
	defer m.db.Close()
	_, err = m.db.QueryContext(m.ctx, "CREATE TABLE IF NOT EXISTS gauge (key varchar(100), value DOUBLE PRECISION);")
	if err != nil {
		return err
	}
	_, err = m.db.QueryContext(m.ctx, "CREATE TABLE IF NOT EXISTS counter (key varchar(100), value int8);")
	if err != nil {
		return err
	}
	return nil

}

func (m *StorageDB) NewGauge(key string, value float64) error {
	return nil
}

func (m *StorageDB) NewCounter(key string, value int64) error {
	return nil
}

func (m *StorageDB) GetCounter(key string) (int64, error) {

	return 1, nil
}

func (m *StorageDB) GetGauge(key string) (float64, error) {

	return 1, nil
}

func (m *StorageDB) Load(file string) error {

	return nil
}

func (m *StorageDB) Save(file string) error {

	return nil

}

func (m *StorageDB) PrintMetric() string {

	return "er"
}
