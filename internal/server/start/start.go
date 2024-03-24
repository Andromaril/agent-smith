package start

import (
	"context"
	"database/sql"

	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/andromaril/agent-smith/internal/serverflag"
)

func Start() (*sql.DB, storage.Storage) {
	var newMetric storage.Storage
	var err error
	var db *sql.DB
	if serverflag.Databaseflag != "" {
		newMetric = &storagedb.StorageDB{Path: serverflag.Databaseflag}
		db, err = newMetric.Init(serverflag.Databaseflag, context.Background())
		if err != nil {
			panic(err)
		}
	} else {
		newMetric = &storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}, WriteSync: serverflag.StoreInterval == 0, Path: serverflag.FileStoragePath}
	}
	return db, newMetric
}
