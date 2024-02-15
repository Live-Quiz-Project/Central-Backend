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
		CreatedAt: 			time.Now(),
		UpdatedAt:			time.Now(),
		DeletedAt: 			gorm.DeletedAt{},
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
		QuizID:					uuid.New(),
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
		QuizID:					uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id","creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String() ,quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

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
		QuizID:					uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id","creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String() ,quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

	// Expected Query
	expectedSQL := "SELECT (.+) FROM \"quiz_history\" WHERE id =.+"
	mock.ExpectQuery(expectedSQL).
		WillReturnRows(row)

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
		QuizID:					id,
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
	row := sqlmock.NewRows([]string{"id", "quiz_id","creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String() ,quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

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
		QuizID:					uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id","creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String() ,quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

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
		QuizID:					uuid.New(),
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
	row := sqlmock.NewRows([]string{"id", "quiz_id","creator_id", "title", "description", "cover_image", "visibility", "time_limit", "have_time_factor", "time_factor", "font_size", "mark", "select_min", "select_max", "case_sensitive"}).
		AddRow(quizH.ID.String(), quizH.QuizID.String() ,quizH.CreatorID.String(), quizH.Title, quizH.Description, quizH.CoverImage, quizH.Visibility, quizH.TimeLimit, quizH.HaveTimeFactor, quizH.TimeFactor, quizH.FontSize, quizH.Mark, quizH.SelectMin, quizH.SelectMax, quizH.CaseSensitive)

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
		QuizID:					uuid.New(),
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

	// Mock Data
	data := &QuestionPool{
		ID:             uuid.New(),
		QuizID:					uuid.New(),
		Order:					1,
		PoolOrder:			1,
		Content:      	"Content",
		Note:						"Note",
		Media:					"Media",
		MediaType:			"MediaType",
		TimeLimit:      20,
		HaveTimeFactor: false,
		TimeFactor:     1,
		FontSize:       24,
	}

	// Expected Query
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"question_pool\" (.+) VALUES (.+)").
		WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Actual Function
	res, err := r.CreateQuestionPool(context.TODO(), db, data)

	// Unit Test
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}