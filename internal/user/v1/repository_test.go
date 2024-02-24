package v1

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestCreateUser(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
	// 	AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"user\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateUser(context.TODO(), data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUsers(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\""
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetUsers(context.TODO())

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetUserByID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.Email).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetUserByEmail(context.TODO(), data.Email)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
	// 	AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"user\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"user\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateUser(context.TODO(), data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
	// 	AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"user\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"user\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteUser(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRestoreUser(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"user\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.RestoreUser(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	//assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestChangePassword(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"user\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.ChangePassword(context.TODO(), data.ID, data.Password)

	// Unit Test
	assert.NoError(t, err)
	//assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetUserByGoogleID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	googleID := "GOOGLE ID"

	data := &User{
		ID:            uuid.New(),
		GoogleId:      &googleID,
		Name:          "Name",
		Email:         "email@email.com",
		Password:      "PASSWORD",
		Image:         "Image",
		DisplayName:   "DisplayName",
		DisplayEmoji:  "DisplayEmoji",
		DisplayColor:  "DisplayColor",
		AccountStatus: "ACTIVE",
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"user\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
		AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"user\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.GoogleId).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetUserByGoogleID(context.TODO(), *data.GoogleId)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateAdmin(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &Admin{
		ID:            uuid.New(),
		Email:         "email@email.com",
		Password:      "PASSWORD",
	}

	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"admin\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "google_id", "name", "email", "password", "image", "display_name", "display_emoji", "display_color", "account_status"}).
	// 	AddRow(data.ID.String(), data.GoogleId, data.Name, data.Email, data.Password, data.Image , data.DisplayName, data.DisplayEmoji, data.DisplayColor, data.AccountStatus)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"user\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"user\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateAdmin(context.TODO(), data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}