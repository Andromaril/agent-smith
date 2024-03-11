package storagedb

import (
	"context"
	"database/sql"

	"github.com/andromaril/agent-smith/internal/server/storage"
)

type StorageDB struct {
	DB   *sql.DB
	Path string
	Ctx  context.Context
}

type Interface interface {
	storage.Storage
	Ping() error
}

func (m *StorageDB) Init(path string, ctx context.Context) (*sql.DB, error) {
	var err error
	m.Ctx = ctx
	m.Path = path
	m.DB, err = sql.Open("pgx", path)
	if err != nil {
		return nil, err
	}
	//defer m.db.Close()
	_, err = m.DB.Exec(`CREATE TABLE IF NOT EXISTS gauge (key varchar(100), value DOUBLE PRECISION)`)
	if err != nil {
		return nil, err
	}
	_, err = m.DB.Exec(`CREATE TABLE IF NOT EXISTS counter (key varchar(100), value int8)`)
	if err != nil {
		return nil, err
	}
	return m.DB, nil

}

func (s *StorageDB) Ping() error {
	return s.DB.Ping()
}

func (m *StorageDB) NewGauge(key string, value float64) error {

	return nil
}

func (m *StorageDB) NewCounter(key string, value int64) error {
	_, err := m.DB.Exec(`
		INSERT INTO counter (key, value) VALUES ($1, $2)`, key, value)
	return err
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

	return "error"
}
