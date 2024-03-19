package storagedb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/andromaril/agent-smith/internal/model"
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
	m.Bootstrap(m.Ctx)
	// //defer m.db.Close()
	// _, err = m.DB.Exec(`CREATE TABLE IF NOT EXISTS gauge (key varchar(100), value DOUBLE PRECISION)`)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = m.DB.Exec(`CREATE TABLE IF NOT EXISTS counter (key varchar(100) UNIQUE NOT NULL, value int8)`)
	// if err != nil {
	// 	return nil, err
	// }
	return m.DB, nil

}

func (m *StorageDB) Bootstrap(ctx context.Context) error {
	// запускаем транзакцию
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gauge (key varchar(100), value DOUBLE PRECISION)`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS counter (key varchar(100) UNIQUE NOT NULL, value int8)`)
	if err != nil {
		return err
	}
	return tx.Commit()
}
func (m *StorageDB) Ping() error {
	return m.DB.Ping()
}

func (m *StorageDB) NewGaugeUpdate(gauge []model.Gauge) error {
	tx, err := m.DB.BeginTx(m.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, value := range gauge {
		_, err = tx.ExecContext(m.Ctx, `
			INSERT INTO gauge (key, value)
			VALUES($1, $2);
		`, value.Key, value.Value)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *StorageDB) NewCounterUpdate(counter []model.Counter) error {
	tx, err := m.DB.BeginTx(m.Ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, value := range counter {
		_, err = tx.ExecContext(m.Ctx, `
			INSERT INTO gauge (key, value)
			VALUES($1, $2);
		`, value.Key, value.Value)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *StorageDB) NewGauge(key string, value float64) error {
	_, err := m.DB.ExecContext(m.Ctx, `
		INSERT INTO gauge (key, value) VALUES ($1, $2)`, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (m *StorageDB) NewCounter(key string, value int64) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO counter (key, value)
	VALUES($1, $2) 
	ON CONFLICT (key) 
	DO UPDATE SET value = counter.value + $2;
`, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (m *StorageDB) GetCounter(key string) (int64, error) {
	var value sql.NullInt64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT value FROM counter WHERE key=$1", key)
	err := rows.Scan(&value)
	if err != nil {
		return 0, err
	}
	if !value.Valid {
		return 0, err
	}
	return value.Int64, nil
}

func (m *StorageDB) GetGauge(key string) (float64, error) {
	var value sql.NullFloat64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT value FROM gauge WHERE key=$1", key)
	err := rows.Scan(&value)
	if err != nil {
		return 0, err
	}
	if !value.Valid {
		return 0, err
	}
	return value.Float64, nil
}

func (m *StorageDB) Load(file string) error {

	return nil
}

func (m *StorageDB) Save(file string) error {

	return nil

}

func (m *StorageDB) GetIntMetric() (map[string]int64, error) {
	counter := make(map[string]int64, 0)
	//gauge := make(map[string]float64, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT key, value FROM counter")
	if err != nil {
		return counter, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var key string
		var value int64
		err = rows.Scan(&key, &value)
		if err != nil {
			return counter, err
		}
		counter[key] = value
	}
	err = rows.Err()
	if err != nil {
		return counter, err
	}
	return counter, nil
}

func (m *StorageDB) GetFloatMetric() (map[string]float64, error) {
	gauge := make(map[string]float64, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT key, value FROM gauge")
	if err != nil {
		return gauge, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var key string
		var value float64
		err = rows.Scan(&key, &value)
		if err != nil {
			return gauge, err
		}
		gauge[key] = value
	}
	err = rows.Err()
	if err != nil {
		return gauge, err
	}
	return gauge, nil
}

func (m *StorageDB) PrintMetric() string {
	counter, err := m.GetIntMetric()
	gauge, err2 := m.GetFloatMetric()
	if err != nil {
		return "error"
	}
	if err2 != nil {
		return "error"
	}
	var result string
	for k1, v1 := range gauge {
		result += fmt.Sprintf("%s: %v\n", k1, v1)
	}
	for k2, v2 := range counter {
		result += fmt.Sprintf("%s: %v\n", k2, v2)
	}
	return result

}
