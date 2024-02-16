package v1

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

func TestTemplate(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.RestoreQuestion(context.TODO(), db, data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
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

	// // Add rows to 'Test' Database

	// Expected Query
	expectedSQL := "UPDATE \"quiz\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuiz(context.TODO(), db, id)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRestoreQuiz(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	quizID := uuid.New()

	// Mock Data
	quiz := &Quiz{
		ID:             quizID,
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
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DeletedAt:      gorm.DeletedAt{},
	}

	mock.ExpectQuery("SELECT (.+) FROM \"quiz\" WHERE \"quiz\".\"id\" = (.+)").
		WithArgs(quizID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive", "created_at", "updated_at", "deleted_at"}).
			AddRow(quiz.ID, quiz.CreatorID, quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive, quiz.CreatedAt, quiz.UpdatedAt, quiz.DeletedAt.Time))

		// Mocking the Update call in the Model chain
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"quiz\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.RestoreQuiz(context.TODO(), db, quizID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateQuizHistory(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	r := NewRepository(db)

	// Variables

	// Mock Data
	quiz := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
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
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"quiz_history\" (.+) VALUES (.+)").
		WithArgs(quiz.ID.String(), quiz.QuizID.String(), quiz.CreatorID.String(), quiz.Title, quiz.Description, quiz.CoverImage, quiz.Visibility, quiz.TimeLimit, quiz.HaveTimeFactor, quiz.TimeFactor, quiz.FontSize, quiz.Mark, quiz.SelectMin, quiz.SelectMax, quiz.CaseSensitive, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuizHistory(context.TODO(), db, quiz)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizHistoryies(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quizH := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String(), quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" "
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(row)

	// Actual Function
	res, err := repo.GetQuizHistories(context.TODO())

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizHistoryByID(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quizH := &QuizHistory{
		ID:             id,
		QuizID:         uuid.New(),
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

	// Add rows in 'Test' database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String(), quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" WHERE id =.+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuizHistoryByID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizHistoriesByQuizID(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quizH := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         id,
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

	// Add rows in 'Test' database
	row := sqlmock.NewRows([]string{"id", "quiz_id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String(), quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" WHERE quiz_id = .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(row)

	// Actual Function
	res, err := repo.GetQuizHistoriesByQuizID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizHistoriesByUserID(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Mock Data
	quizH := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String(), quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" WHERE creator_id = .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(row)

	// Actual Function
	res, err := repo.GetQuizHistoriesByUserID(context.TODO(), id)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuizHistoriesByQuizIDAndCreatedDate(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables

	// Mock Data
	quizH := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
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

	// Add rows in 'Test' database
	row := sqlmock.NewRows([]string{"id", "quiz_id", "creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String(), quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" WHERE .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(row)

	// Actual Function
	res, err := repo.GetQuizHistoryByQuizIDAndCreatedDate(context.TODO(), quizH.QuizID, quizH.CreatedAt)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuizHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables

	// Mock Data
	quizH := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
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

	// Expected Query
	expectedSQL := "UPDATE \"quiz_history\" SET (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuizHistory(context.TODO(), db, quizH)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuizHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Variables
	id := uuid.New()

	// Expected Query
	expectedSQL := "UPDATE \"quiz_history\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuizHistory(context.TODO(), db, id)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateQuestionPool(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	r := NewRepository(db)

	// Expected Query
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"question_pool\" (.+) VALUES (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuestionPool(context.TODO(), db, &QuestionPool{})

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionPoolByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "order", "pool_order", "content", "note", "media", "media_type", "time_limit", "have_time_factor", "time_factor", "font_size"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.Order, data.PoolOrder, data.Content, data.Note, data.Media, data.MediaType, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_pool\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID.String()).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionPoolByID(context.TODO(), data.ID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionPoolsByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "order", "pool_order", "content", "note", "media", "media_type", "time_limit", "have_time_factor", "time_factor", "font_size"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.Order, data.PoolOrder, data.Content, data.Note, data.Media, data.MediaType, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_pool\" WHERE quiz_id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID.String()).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionPoolsByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetDeleteQuestionPoolsByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "order", "pool_order", "content", "note", "media", "media_type", "time_limit", "have_time_factor", "time_factor", "font_size"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.Order, data.PoolOrder, data.Content, data.Note, data.Media, data.MediaType, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_pool\" WHERE quiz_id =(.+)"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID.String()).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetDeleteQuestionPoolsByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestionPool(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question_pool\" SET (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuestionPool(context.TODO(), db, data)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestionPool(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question_pool\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuestionPool(context.TODO(), db, data.ID)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRestoreQuestionPool(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "order", "pool_order", "content", "note", "media", "media_type", "time_limit", "have_time_factor", "time_factor", "font_size"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.Order, data.PoolOrder, data.Content, data.Note, data.Media, data.MediaType, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize)

	expectedSQL := "SELECT (.+) FROM \"question_pool\" WHERE .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"question_pool\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.RestoreQuestionPool(context.TODO(), db, uuid.New())

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateQuestionPoolHistory(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	r := NewRepository(db)

	// Expected Query
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"question_pool_history\" (.+) VALUES (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuestionPoolHistory(context.TODO(), db, &QuestionPoolHistory{})

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionPoolHistoriesByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPoolHistory{
		ID:             uuid.New(),
		QuestionPoolID: uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "question_pool_id", "quiz_id", "order", "pool_order", "content", "note", "media", "media_type", "time_limit", "have_time_factor", "time_factor", "font_size"}).
		AddRow(data.ID.String(), data.QuestionPoolID.String(), data.QuizID.String(), data.Order, data.PoolOrder, data.Content, data.Note, data.Media, data.MediaType, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_pool_history\" WHERE .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID.String()).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionPoolHistoriesByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestionPoolHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPoolHistory{
		ID:             uuid.New(),
		QuestionPoolID: uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question_pool_history\" SET (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuestionPoolHistory(context.TODO(), db, data)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestionPoolHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		Order:          1,
		PoolOrder:      -1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question_pool\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuestionPool(context.TODO(), db, data.ID)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateQuestion(t *testing.T) {
	// Test Setup
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	r := NewRepository(db)

	id := uuid.New()
	// Mock Data
	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}
	
	// Expected Query
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"question\" (.+) VALUES (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuestion(context.TODO(), db, data)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestions(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	// Mock Data
	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\""
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestions(context.TODO())

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()

	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\" WHERE id = .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionByID(context.TODO(), data.ID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionsByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()

	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\" WHERE quiz_id = .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionsByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetDeleteQuestionsByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()

	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\" WHERE quiz_id = .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetDeleteQuestionsByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionByQuizIDAndOrder(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	// Mock Data
	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Add rows to 'Test' Database
	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\" WHERE .+ "
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID,data.Order).
		WillReturnRows(sample)

	// Actual Function
	res, err := repo.GetQuestionByQuizIDAndOrder(context.TODO(), data.QuizID, data.Order)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestion(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()

	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question\" SET (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuestion(context.TODO(), db, data)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestion(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// Expected Query
	expectedSQL := "UPDATE \"question\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuestion(context.TODO(), db, data.ID)

	// Unit Test
	assert.Nil(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRestoreQuestion(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &Question{
		ID:             uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID).
		WillReturnRows(sample)

		// Mocking the Update call in the Model chain
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"question\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.RestoreQuestion(context.TODO(), db, data.QuizID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateQuestionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}


	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"question_history\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
	// 	AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"question\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateQuestionHistory(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionHistories(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetQuestionHistories(context.TODO())

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionHistoryByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id","quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetQuestionHistoryByID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionHistoriesByQuizID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuizID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetQuestionHistoriesByQuizID(context.TODO(), data.QuizID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionHistoriesByQuestionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetQuestionHistoriesByQuestionID(context.TODO(), data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetQuestionHistoryByQuestionIDAndCreatedDate(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"question_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID, data.CreatedAt).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"question\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetQuestionHistoryByQuestionIDAndCreatedDate(context.TODO(), data.QuestionID, data.CreatedAt)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
	// 	AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"question\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"question_history\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateQuestionHistory(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	id := uuid.New()
	data := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		QuizID:         uuid.New(),
		QuestionPoolID: &id,
		PoolOrder:      -1,
		PoolRequired:   false,
		Type:           "Type",
		Order:          1,
		Content:        "Content",
		Note:           "Note",
		Media:          "Media",
		MediaType:      "MediaType",
		UseTemplate:    false,
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       16,
		LayoutIdx:      1,
		SelectMin:      1,
		SelectMax:      4,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"quiz\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "quiz_id", "question_pool_id", "pool_order","pool_required", "type", "order", "content", "note", "media", "media_type", "use_template","time_limit", "have_time_factor", "time_factor", "font_size", "layout_idx", "select_min", "select_max"}).
	// 	AddRow(data.ID.String(), data.QuizID.String(), data.QuestionPoolID.String(), data.PoolOrder, data.PoolRequired,data.Type, data.Order, data.Content, data.Note, data.Media, data.MediaType, data.UseTemplate, data.TimeLimit, data.HaveTimeFactor, data.TimeFactor, data.FontSize, data.LayoutIdx, data.SelectMin, data.SelectMax)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"question\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"question_history\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteQuestionHistory(context.TODO(), db , data.ID)

	// Unit Test
	assert.NoError(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateChoiceOption(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.QuestionID.String(), data.order, data.content, data.mark, data.color, data.correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateChoiceOption(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceOptionByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceOptionByID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceOptionsByQuestionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceOptionsByQuestionID(context.TODO(), data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetDeleteChoiceOptionsByQuestionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetDeleteChoiceOptionsByQuestionID(context.TODO(), data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceAnswersByQuestionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID, data.Correct).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceAnswersByQuestionID(context.TODO(), data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateChoiceOption(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.QuestionID.String(), data.order, data.content, data.mark, data.color, data.correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"option_choice\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateChoiceOption(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteChoiceOption(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.QuestionID.String(), data.order, data.content, data.mark, data.color, data.correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"option_choice\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteChoiceOption(context.TODO(), db, data.ID)

	// Unit Test
	assert.NoError(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRestoreChoiceOption(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOption{
		ID:             uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"option_choice\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.RestoreChoiceOption(context.TODO(), db, data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateChoiceOptionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQL).
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()). // Number of Data in Struct
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.CreateChoiceOptionHistory(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceOptionHistoryByID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.ID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceOptionHistoryByID(context.TODO(), data.ID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetOptionChoiceHistories(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetOptionChoiceHistories(context.TODO())

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceOptionHistoriesByQuestionID(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceOptionHistoriesByQuestionID(context.TODO(), data.QuestionID)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetChoiceOptionHistoryByQuestionIDAndContent(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
		AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	mock.ExpectQuery(expectedSQL).
		WithArgs(data.QuestionID, data.Content).
		WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	// mock.ExpectBegin()
	// mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()

	// Actual Function
	res, err := repo.GetChoiceOptionHistoryByQuestionIDAndContent(context.TODO(), data.QuestionID , data.Content)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateChoiceOptionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := repo.UpdateChoiceOptionHistory(context.TODO(), db, data)

	// Unit Test
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteChoiceOptionHistory(t *testing.T) {
	// Setup Test
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	// Mock Data
	data := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: uuid.New(),
		QuestionID:			uuid.New(),
		Order:					1,
		Content:        "Content",
		Mark:						10,
		Color:					"WHITE",
		Correct:				true,
	}

	// ===== CREATE  =====
	// expectedSQL := "INSERT INTO \"option_choice_history\" (.+) VALUES (.+)"
	// mock.ExpectBegin()
	// mock.ExpectExec(expectedSQL).
	// 	WithArgs(sqlmock.AnyArg()). // Number of Data in Struct
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// mock.ExpectCommit()


	// ===== GET RESTORE =====
	// sample := sqlmock.NewRows([]string{"id", "option_choice_id", "question_id", "order", "content","mark", "color", "correct"}).
	// 	AddRow(data.ID.String(), data.ChoiceOptionID.String(), data.QuestionID.String(), data.Order, data.Content, data.Mark, data.Color, data.Correct)

	// // Expected Query
	// expectedSQL := "SELECT (.+) FROM \"option_choice_history\" .+"
	// mock.ExpectQuery(expectedSQL).
	// 	WithArgs(data.QuizID).
	// 	WillReturnRows(sample)

	// ===== UPDATE DELETE RESTORE =====
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"option_choice_history\" SET .+").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	err := repo.DeleteChoiceOptionHistory(context.TODO(), db, data.ID)

	// Unit Test
	assert.NoError(t, err)
	// assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}