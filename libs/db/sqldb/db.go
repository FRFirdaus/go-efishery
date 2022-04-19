package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	// when need use an other driver
	// you can register like this below on your own module
	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

var store sync.Map

// open is internali open database source
func open(config *Configuration) (*sql.DB, error) {
	db, err := dburl.Open(config.Dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// ensure connection establish
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}

// Close closes all connections
// If passing key will close connection of specified keys
func Close(dbKey ...string) error {
	if len(dbKey) > 1 {
		value, ok := store.Load(dbKey[0])
		if !ok {
			return fmt.Errorf("Invalid key")
		}

		if err := (value.(*sql.DB)).Close(); err != nil {
			return err
		}

		store.Delete(dbKey[0])
		return nil
	}

	store.Range(func(key, value interface{}) bool {
		log.Printf("Closing SQL DB: %s", key)

		if err := (value.(*sql.DB)).Close(); err != nil {
			log.Println(err)
		}

		store.Delete(key)

		return true
	})

	return nil
}

// Open is open sql connection with passing configuration by default using "default" key
// Optional can passing specific key
func Open(config *Configuration) error {
	if config.Key == "" {
		return fmt.Errorf("Key db is required")
	}

	// validate dsn and parse additional option
	// do silent error
	additionalConfiguration, err := ParseURL(config.Dsn)
	if err != nil {
		return err
	}

	db, err := open(config)
	if err != nil {
		return err
	}

	// set additional configuration
	db.SetConnMaxLifetime(time.Duration(additionalConfiguration.MaxConnLifeTime))
	db.SetMaxIdleConns(additionalConfiguration.MaxIdleConnections)
	db.SetMaxOpenConns(additionalConfiguration.MaxOpenConnections)

	store.Store(config.Key, db)

	return nil
}

// DB is return connection of database
func DB(key string) *sql.DB {
	db, ok := store.Load(key)
	if !ok {
		log.Fatal("Please open db first", key)
	}

	return db.(*sql.DB)
}
