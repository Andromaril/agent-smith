package storagedb

import (
	"context"
	"database/sql"
)

// type Store struct {
// 	// Поле conn содержит объект соединения с СУБД
// 	conn *sql.DB
// }

// // NewStore возвращает новый экземпляр PostgreSQL-хранилища
// func NewStore(conn *sql.DB) *Store {
// 	return &Store{conn: conn}
// }

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
	var err error
	var count int
	row := m.db.QueryRowContext(m.ctx, "SELECT COUNT(*) FROM counter WHERE key=$1", key)
	if err = row.Err(); err != nil {
		return err
	}
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = m.db.ExecContext(m.ctx, "INSERT INTO counter (key, value) VALUES ($1, $2)", key, value)
		if err != nil {
			return err
		}
	} else {
		var counterVal int64
		row2 := m.db.QueryRowContext(m.ctx, "SELECT value FROM counter WHERE key=$1", key)
		if err = row2.Err(); err != nil {
			return err
		}
		err = row2.Scan(&counterVal)
		if err != nil {
			return err
		}
		_, err = m.db.ExecContext(m.ctx, "UPDATE counter SET value=$1 WHERE key=$2", counterVal+value, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *StorageDB) GetCounter(key string) (int64, error) {

	var v sql.NullInt64
	row := m.db.QueryRowContext(m.ctx, "SELECT value FROM counter WHERE key=$1", key)
	if err := row.Err(); err != nil {
		return 0, nil
	}
	if err := row.Scan(&v); err != nil {
		return 0, nil
	}
	if !v.Valid {
		return 0, nil
	}
	return v.Int64, nil
}

func (m *StorageDB) GetGauge(key string) (float64, error) {

	var v sql.NullFloat64
	row := m.db.QueryRowContext(m.ctx, "SELECT value FROM gauge WHERE key=$1", key)
	if err := row.Err(); err != nil {
		return 0, nil
	}
	if err := row.Scan(&v); err != nil {
		return 0, nil
	}
	if !v.Valid {
		return 0, nil
	}
	return v.Float64, nil
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
