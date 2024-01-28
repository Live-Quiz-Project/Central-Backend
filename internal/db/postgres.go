package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (*Database, error) {
	
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	er := sqlDB.Ping()
	if er != nil {
		return nil, er
	}

	return &Database{db: db}, nil
}

func NewTestDatabase() (*Database, error) {
	
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := "test_database"

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	er := sqlDB.Ping()
	if er != nil {
		return nil, er
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	sqlDB, err := d.db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.Close()
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
