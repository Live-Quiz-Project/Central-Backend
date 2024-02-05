package v1

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Create repository instance with the test database
	repo := NewRepository(testDB)

	// Create a user to be inserted
	user := &User{
		ID:            uuid.New(),
		Name:          "Beta Tester",
		Email:         "test-user@example.com",
		Password:      "hashed_password",
		Image:         "profile.jpg",
		DisplayName:   "Test",
		DisplayEmoji:  ":smile:",
		DisplayColor:  "#3498db",
		AccountStatus: "active",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		// DeletedAt:   gorm.DeletedAt{}, // You may skip this field if it's not set
	}

	createdUser, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)

	assert.NotNil(t, createdUser)

	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Image, createdUser.Image)
	assert.Equal(t, user.DisplayName, createdUser.DisplayName)
	assert.Equal(t, user.DisplayEmoji, createdUser.DisplayEmoji)
	assert.Equal(t, user.DisplayColor, createdUser.DisplayColor)

	// You can also retrieve the user from the database to make sure it was actually inserted
	var retrievedUser User
	testDB.First(&retrievedUser, "id = ?", user.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Image, retrievedUser.Image)
	assert.Equal(t, user.DisplayName, retrievedUser.DisplayName)
	assert.Equal(t, user.DisplayEmoji, retrievedUser.DisplayEmoji)
	assert.Equal(t, user.DisplayColor, retrievedUser.DisplayColor)
}
