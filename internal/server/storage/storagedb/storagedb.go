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
	var err error
	var count int
	gaugeCountQuery := m.db.QueryRowContext(m.ctx, "SELECT COUNT(*) FROM gauge WHERE key=$1", key)
	if err = gaugeCountQuery.Err(); err != nil {
		return err
	}
	err = gaugeCountQuery.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = m.db.ExecContext(m.ctx, "INSERT INTO gauge (key, value) VALUES ($1, $2)", key, value)
		if err != nil {
			return err
		}
	} else {
		_, err = m.db.ExecContext(m.ctx, "UPDATE gauge SET value=$1 WHERE key=$2", value, key)
		if err != nil {
			return err
		}
	}
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
