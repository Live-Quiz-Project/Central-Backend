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

func TestCreateQuiz(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	r := NewRepository(db)

	// Variables

	// Mock Data
	quiz := &Quiz{
		ID:             uuid.New(),
		CreatorID:      uuid.New(),
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows in 'Test' Database

	// Expected Query
	expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(quiz.ID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuiz(context.TODO(), db, quiz)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizzesByUserID(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quiz := &Quiz{
		ID:             uuid.New(),
		CreatorID:      id,
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows in 'Test' database
	quizzes := sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quiz.ID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz\" WHERE creator_id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(id.String()).
		WillReturnRows(quizzes)

	// Actual Function
	res, err := repo.GetQuizzesByUserID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quiz := &Quiz{
		ID:             id,
		CreatorID:      uuid.New(),
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows to 'Test' Database
	quizzes := sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quiz.ID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(id.String()).
		WillReturnRows(quizzes)

	// Actual Function
	res, err := repo.GetQuizByID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetDeleteQuizByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quiz := &Quiz{
		ID:             id,
		CreatorID:      uuid.New(),
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows to 'Test' Database
	quizzes := sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quiz.ID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(id.String()).
		WillReturnRows(quizzes)

	// Actual Function
	res, err := repo.GetDeleteQuizByID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuiz(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables

	// Mock Data
	quiz := &Quiz{
		ID:             uuid.New(),
		CreatorID:      uuid.New(),
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows to 'Test' Database
	sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quiz.ID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive)

	// Expected Query
	expectedSQL := "UPDATE \"quiz\" SET (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuiz(context.TODO(), db, quiz)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuiz(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quiz := &Quiz{
		ID:             id,
		CreatorID:      uuid.New(),
		Title:          "Test Title",
		Description:    "Test Description",
		CoverImage:     "Test CoverImage",
		Visibility:     "PRIVATE",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
		Mark:           10,
		SelectMin:      1,
		SelectMax:      1,
		CaseSensitive:  true,
	}

	// Add rows to 'Test' Database
	sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(id.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive)

	// Expected Query
	expectedSQL := "UPDATE \"quiz\" SET \"deleted_at\"=$1 WHERE \"quiz\".\"id\" = $2 AND \"quiz\".\"deleted_at\" IS NULL"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
	  WithArgs("2024-02-13T16:58:58.390161Z",id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuiz(context.TODO(), db, id)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}
