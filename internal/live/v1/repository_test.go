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

func NewTestRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

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

func TestGetLiveQuizSessionBySessionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTestRepository(db)

	// Mock Data
	data := &Session{
		ID:                   uuid.New(),
		HostID:               uuid.New(),
		QuizID:               uuid.New(),
		Status:               "Status",
		ExemptedQuestionIDs: 	nil,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"live_quiz_session\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "host_id", "quiz_id","status", "exempted_question_ids"}).
		AddRow(data.ID.String(), data.HostID.String(), data.QuizID.String(), data.Status, data.ExemptedQuestionIDs)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"live_quiz_session\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"live_quiz_session\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetLiveQuizSessionBySessionID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetLiveQuizSessionsByUserID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTestRepository(db)

	// Mock Data
	data := &Session{
		ID:                   uuid.New(),
		HostID:               uuid.New(),
		QuizID:               uuid.New(),
		Status:               "Status",
		ExemptedQuestionIDs: 	nil,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"live_quiz_session\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "host_id", "quiz_id","status", "exempted_question_ids"}).
		AddRow(data.ID.String(), data.HostID.String(), data.QuizID.String(), data.Status, data.ExemptedQuestionIDs)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"live_quiz_session\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.HostID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"live_quiz_session\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetLiveQuizSessionsByUserID(context.TODO(), data.HostID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateLiveQuizSession(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTestRepository(db)

	// Mock Data
	data := &Session{
		ID:                   uuid.New(),
		HostID:               uuid.New(),
		QuizID:               uuid.New(),
		Status:               "Status",
		ExemptedQuestionIDs: 	nil,
	}

	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"live_quiz_session\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "host_id", "quiz_id","status", "exempted_question_ids"}).
	// 	AddRow(data.ID.String(), data.HostID.String(), data.QuizID.String(), data.Status, data.ExemptedQuestionIDs)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"live_quiz_session\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.ID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"live_quiz_session\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateLiveQuizSession(context.TODO(), data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetLiveQuizSessions(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTestRepository(db)

	// Mock Data
	data := &Session{
		ID:                   uuid.New(),
		HostID:               uuid.New(),
		QuizID:               uuid.New(),
		Status:               "Status",
		ExemptedQuestionIDs: 	nil,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"live_quiz_session\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "host_id", "quiz_id","status", "exempted_question_ids"}).
		AddRow(data.ID.String(), data.HostID.String(), data.QuizID.String(), data.Status, data.ExemptedQuestionIDs)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"live_quiz_session\""
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"live_quiz_session\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetLiveQuizSessions(context.TODO())

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

// Test Live Quiz Session 
