package v1

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {

	log.Println("Start Test on User")
	setup()

	exitCode := m.Run()

	teardown()

	os.Exit(exitCode)
}

func setup() {
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "DBuser"
	dbPass := "DBpass"
	dbName := "test_database"

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlFile, err := os.ReadFile("../../db/init_table.sql")
	if err != nil {
		panic(err)
	}

	db.Exec(string(sqlFile))

	testDB = db
}

func teardown() {
	// Add logic to delete records or perform any other cleanup in the database
	// For example, you might use testDB.Exec("DELETE FROM users") to delete all user records
	testDB.Exec("TRUNCATE TABLE \"user\" RESTART IDENTITY CASCADE;")
}
