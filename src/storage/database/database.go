package database

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type Driver string

const (
	Postgres Driver = "POSTRGRES"
	SQLite   Driver = "SQLITE"
)

func New(driver Driver) {
	switch driver {
	case SQLite:
		newSQLiteDB()
	}
}

func newSQLiteDB() {
	once.Do(func() {
		var err error

		db, err = gorm.Open(sqlite.Open("db.sqlite"))
		if err != nil {
			log.Fatalf("Can't connect to the database %v", err)
		}

		fmt.Println("✅ Connected on sqlite successfully")
	})
}

func DB() *gorm.DB {
	return db
}
