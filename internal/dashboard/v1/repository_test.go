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

func TestGetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &AnswerResponse{
		ID:                uuid.New(),
		LiveQuizSessionID: uuid.New(),
		ParticipantID:     uuid.New(),
		Type:              "Type",
		QuestionID:        uuid.New(),
		Answer:            "Answer",
		UseTime:           5,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"answer_response\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "live_quiz_session_id", "participant_id", "type", "question_id", "answer", "use_time"}).
		AddRow(data.ID.String(), data.LiveQuizSessionID.String(), data.ParticipantID.String(), data.Type, data.QuestionID.String(), data.Answer, data.UseTime)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"answer_response\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.LiveQuizSessionID, data.QuestionID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"answer_response\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(context.TODO(), data.LiveQuizSessionID, data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetAnswerResponsesByLiveQuizSessionIDAndParticipantID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &AnswerResponse{
		ID:                uuid.New(),
		LiveQuizSessionID: uuid.New(),
		ParticipantID:     uuid.New(),
		Type:              "Type",
		QuestionID:        uuid.New(),
		Answer:            "Answer",
		UseTime:           5,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"answer_response\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "live_quiz_session_id", "participant_id", "type", "question_id", "answer", "use_time"}).
		AddRow(data.ID.String(), data.LiveQuizSessionID.String(), data.ParticipantID.String(), data.Type, data.QuestionID.String(), data.Answer, data.UseTime)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"answer_response\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.LiveQuizSessionID, data.ParticipantID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"answer_response\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetAnswerResponsesByLiveQuizSessionIDAndParticipantID(context.TODO(), data.LiveQuizSessionID, data.ParticipantID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetParticipantByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	userID := uuid.New()

	data := &Participant{
		ID:                uuid.New(),
		UserID:            &userID,
		LiveQuizSessionID: uuid.New(),
		Status:            "ACTIVE",
		Name:              "Name",
		Marks:             100,
	}

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "user_id", "live_quiz_session_id", "status", "name", "marks"}).
		AddRow(data.ID.String(), data.UserID.String(), data.LiveQuizSessionID.String(), data.Status, data.Name, data.Marks)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"participant\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetParticipantByID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetOrderParticipantsByLiveQuizSessionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	userID := uuid.New()

	data := &Participant{
		ID:                uuid.New(),
		UserID:            &userID,
		LiveQuizSessionID: uuid.New(),
		Status:            "ACTIVE",
		Name:              "Name",
		Marks:             100,
	}

	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "user_id", "live_quiz_session_id", "status", "name", "marks"}).
		AddRow(data.ID.String(), data.UserID.String(), data.LiveQuizSessionID.String(), data.Status, data.Name, data.Marks)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"participant\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.LiveQuizSessionID).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetOrderParticipantsByLiveQuizSessionID(context.TODO(), data.LiveQuizSessionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}
