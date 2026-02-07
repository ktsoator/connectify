package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	Email     string `gorm:"unique"`
	Password  string
	Nickname  string
	Intro     string
	CreatedAt int64
	UpdatedAt int64
}

var (
	// ErrDuplicateEmail is returned when the email already exists in the database
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrRecordNotFound is returned when a record is not found in the database
	ErrRecordNotFound = errors.New("record not found")
)

const (
	mysqlDuplicateEntryErrCode uint16 = 1062
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (u *UserDAO) Insert(ctx context.Context, user UserModel) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := u.db.WithContext(ctx).Create(&user).Error

	if err != nil {
		// Use errors.As to check if the error is a MySQL driver error.
		// It unwraps the error if it was wrapped by other layers (like GORM).
		var mysqlErr *mysql.MySQLError
		// 1062 is the MySQL error code for "Duplicate entry"
		// This happens when a unique constraint (like the email) is violated.
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrCode {
			return ErrDuplicateEmail
		}
		// If it's not a duplicate email error, return the original error (e.g., db connection lost)
		// We must return the error so the caller knows something went wrong.
		return err
	}
	return nil
}

func (u *UserDAO) FindByEmail(ctx context.Context, email string) (UserModel, error) {
	var user UserModel
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserModel{}, ErrRecordNotFound
		}
		return UserModel{}, err
	}
	return user, nil
}

func (u *UserDAO) FindByID(ctx context.Context, id int64) (UserModel, error) {
	var user UserModel
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserModel{}, ErrRecordNotFound
		}
		return UserModel{}, err
	}
	return user, nil
}

func (u *UserDAO) UpdateById(ctx context.Context, user UserModel) error {
	return u.db.WithContext(ctx).Model(&user).Where("id = ?", user.ID).
		Updates(map[string]any{
			"nickname":   user.Nickname,
			"intro":      user.Intro,
			"updated_at": time.Now().UnixMilli(),
		}).Error
}
