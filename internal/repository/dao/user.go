package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (u *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	return u.db.Create(&user).Error
}

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt int64
	UpdatedAt int64
}
