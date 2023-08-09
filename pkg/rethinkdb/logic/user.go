package logic

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	UserName       string         `json:"user_name"`
	HashedPassword string         `json:"hashed_password"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func NewUser(userName string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		UserName:       userName,
		HashedPassword: string(hashedPassword),
	}

	return &user, nil
}

func (user *User) Create(db *gorm.DB) error {
	user.ID = NewULID().String()
	result := db.Create(user)
	return result.Error

}

func GetUserByUserName(db *gorm.DB, userName string) (*User, error) {
	var user User

	result := db.Where("user_name = ?", userName).First(&user)

	return &user, result.Error
}
