package v1

import (
	"context"
	"testing"
	"time"

	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Create repository instance with the test database
	repo := NewRepository(testDB)
	user_repo := u.NewRepository(testDB)

	tx := testDB.Begin()
	defer tx.Rollback()

	// Create a user to be inserted
	user := &u.User{
		ID:            uuid.New(),
		Name:          "QuizTest User",
		Email:         "quiz-user@tester.com",
		Password:      "hashed_password",
		Image:         "profile.jpg",
		DisplayName:   "Quiz Test",
		DisplayEmoji:  ":smile:",
		DisplayColor:  "#3498db",
		AccountStatus: "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdUser, err := user_repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)

	assert.NotNil(t, createdUser)

	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Image, createdUser.Image)
	assert.Equal(t, user.DisplayName, createdUser.DisplayName)
	assert.Equal(t, user.DisplayEmoji, createdUser.DisplayEmoji)
	assert.Equal(t, user.DisplayColor, createdUser.DisplayColor)

	quiz := &Quiz{
		ID:             uuid.New(),
		CreatorID:      user.ID,
		Title:          "Test Create Quiz",
		Description:    "TestCreateQuiz Function",
		CoverImage:     "cover_image.jpg",
		Visibility:     "public",
		TimeLimit:      60,
		HaveTimeFactor: true,
		TimeFactor:     2,
		FontSize:       16,
		Mark:           10,
		SelectUpTo:     5,
		CaseSensitive:  false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdQuiz, err := repo.CreateQuiz(context.Background(), tx, quiz)
	tx.Commit()

	// You can also retrieve the user from the database to make sure it was actually inserted
	retrievedQuiz, err := repo.GetQuizByID(context.Background(), quiz.ID)

	assert.Equal(t, createdQuiz.ID, retrievedQuiz.ID)
	assert.Equal(t, createdQuiz.CreatorID, retrievedQuiz.CreatorID)
	assert.Equal(t, createdQuiz.Title, retrievedQuiz.Title)
	assert.Equal(t, createdQuiz.Description, retrievedQuiz.Description)
	assert.Equal(t, createdQuiz.CoverImage, retrievedQuiz.CoverImage)
	assert.Equal(t, createdQuiz.Visibility, retrievedQuiz.Visibility)
}
