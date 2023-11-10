package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ---------- User related repository methods ---------- //
func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	res := r.db.WithContext(ctx).Create(user)
	if res.Error != nil {
		return &User{}, res.Error
	}
	return user, nil
}

func (r *repository) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	res := r.db.WithContext(ctx).Unscoped().Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	res := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	res := r.db.WithContext(ctx).Save(user)
	if res.Error != nil {
		return &User{}, res.Error
	}
	return user, nil
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repository) RestoreUser(ctx context.Context, id uuid.UUID) error {
	var user User
	res := r.db.WithContext(ctx).Unscoped().First(&user, id)
	if res.Error != nil {
		return res.Error
	}
	res = r.db.WithContext(ctx).Model(&user).Updates(User{DeletedAt: gorm.DeletedAt{}})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// ---------- Admin related repository methods ---------- //
func (r *repository) CreateAdmin(ctx context.Context, admin *Admin) (*Admin, error) {
	res := r.db.WithContext(ctx).Create(admin)
	if res.Error != nil {
		return &Admin{}, res.Error
	}
	return admin, nil
}
