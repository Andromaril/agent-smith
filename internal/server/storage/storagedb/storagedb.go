package storagedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	var pgErr *pgconn.PgError
	tries := 0
	wait := 1
	for {
		if tries > 3 {
			log.Println("Can't connect to DB")
			e := errormetric.NewMetricError(err)
			return nil, fmt.Errorf("fatal start a transaction %q", e.Error())
		}
		m.DB, err = sql.Open("pgx", path)
		if err != nil {
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				log.Println("Connection error. Trying to reconnect...")
				time.Sleep(time.Duration(wait) * time.Second)
				wait += 2
				tries++
				continue
			}
			log.Println("connect DB critical error", err)
		}
		break
	}
	m.Bootstrap(m.Ctx)
	return m.DB, nil

}

func (m *StorageDB) Bootstrap(ctx context.Context) error {
	// запускаем транзакцию
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()
	tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS gauge (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value DOUBLE PRECISION NOT NULL
		);
	`)
	tx.ExecContext(m.Ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value int8
		);
	`)
	return tx.Commit()
}

func (m *StorageDB) Ping() error {
	return m.DB.Ping()
}

func (m *StorageDB) NewGaugeUpdate(gauge []model.Gauge) error {
	tx, err := m.DB.BeginTx(m.Ctx, nil)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	defer tx.Rollback()
	for _, value := range gauge {
		_, err = tx.ExecContext(m.Ctx, `
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`, value.Key, value.Value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				e := errormetric.NewMetricError(err)
				return fmt.Errorf("error insert %q", e.Error())
			}
		}
	}

	tx.Commit()
	return nil
}

func (m *StorageDB) NewCounterUpdate(counter []model.Counter) error {
	tx, err := m.DB.BeginTx(m.Ctx, nil)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("fatal start a transaction %q", e.Error())
	}
	defer tx.Rollback()
	for _, value := range counter {
		_, err = tx.ExecContext(m.Ctx, `
			INSERT INTO counter (key, value)
			VALUES($1, $2)
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`, value.Key, value.Value)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				e := errormetric.NewMetricError(err)
				return fmt.Errorf("error insert %q", e.Error())
			}
		}
	}
	tx.Commit()
	return nil
}

func (m *StorageDB) NewGauge(key string, value float64) error {
	_, err := m.DB.ExecContext(m.Ctx, `
	INSERT INTO gauge (key, value)
	VALUES($1, $2) 
	ON CONFLICT (key) 
	DO UPDATE SET value = $2;`, key, value)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error insert %q", e.Error())
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
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error insert %q", e.Error())
	}
	return nil
}

func (m *StorageDB) GetCounter(key string) (int64, error) {
	var value sql.NullInt64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT value FROM counter WHERE key=$1", key)
	err := rows.Scan(&value)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return 0, fmt.Errorf("error select %q", e.Error())
	}
	if !value.Valid {
		e := errormetric.NewMetricError(err)
		return 0, fmt.Errorf("invalid value %q", e.Error())
	}
	return value.Int64, nil
}

func (m *StorageDB) GetGauge(key string) (float64, error) {
	var value sql.NullFloat64
	rows := m.DB.QueryRowContext(m.Ctx, "SELECT value FROM gauge WHERE key=$1", key)
	err := rows.Scan(&value)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return 0, fmt.Errorf("error select %q", e.Error())
	}
	if !value.Valid {
		e := errormetric.NewMetricError(err)
		return 0, fmt.Errorf("invalid value %q", e.Error())
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
		e := errormetric.NewMetricError(err)
		return counter, fmt.Errorf("error select %q", e.Error())
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var key string
		var value int64
		err = rows.Scan(&key, &value)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return counter, fmt.Errorf("not int64 metric %q", e.Error())
		}
		counter[key] = value
	}
	err = rows.Err()
	if err != nil {
		e := errormetric.NewMetricError(err)
		return counter, fmt.Errorf("error %q", e.Error())
	}
	return counter, nil
}

func (m *StorageDB) GetFloatMetric() (map[string]float64, error) {
	gauge := make(map[string]float64, 0)
	rows, err := m.DB.QueryContext(m.Ctx, "SELECT key, value FROM gauge")
	if err != nil {
		e := errormetric.NewMetricError(err)
		return gauge, fmt.Errorf("error select %q", e.Error())
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var key string
		var value float64
		err = rows.Scan(&key, &value)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return gauge, fmt.Errorf("not float64 metric %q", e.Error())
		}
		gauge[key] = value
	}
	err = rows.Err()
	if err != nil {
		e := errormetric.NewMetricError(err)
		return gauge, fmt.Errorf("error %q", e.Error())
	}
	return gauge, nil
}

func (m *StorageDB) PrintMetric() string {
	counter, err := m.GetIntMetric()
	gauge, err2 := m.GetFloatMetric()
	if err != nil {
		e := errormetric.NewMetricError(err)
		return e.Error()
	}
	if err2 != nil {
		e := errormetric.NewMetricError(err)
		return e.Error()
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
