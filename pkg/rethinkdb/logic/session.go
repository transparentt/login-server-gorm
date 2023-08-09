package logic

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Session struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	UserULID    string         `json:"user_ul_id"`
	AccessToken string         `json:"access_token"`
	Expired     time.Time      `json:"expired"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewSession(userULID string, accessToken string, expired time.Time) Session {
	return Session{
		UserULID:    userULID,
		AccessToken: accessToken,
		Expired:     expired,
	}
}

func (s *Session) Create(db *gorm.DB) error {
	s.ID = NewULID().String()
	result := db.Create(s)
	return result.Error

}

func GetSessionByUserULID(db *gorm.DB, userULID string) (*Session, error) {

	var session Session

	result := db.Where("user_ul_id = ?", userULID).First(&session)

	return &session, result.Error
}

func UpdateSession(db *gorm.DB, session Session) (*Session, error) {

	result := db.Save(&session)

	return &session, result.Error
}

func CheckSession(db *gorm.DB, userULID string, accessToken string) (*Session, error) {
	session, err := GetSessionByUserULID(db, userULID)
	if session.ID == "" {
		return nil, errors.New("no session")
	}

	if session.AccessToken != accessToken {
		return nil, errors.New("wrong access token")
	}

	if session.Expired.Before(time.Now()) {
		return nil, errors.New("expired access token")
	}

	// Update session with the new access token and expired date.
	newAccessToken, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newExpired := time.Now().Add(time.Hour * 6)

	session.AccessToken = string(newAccessToken)
	session.Expired = newExpired

	_, err = UpdateSession(db, *session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

type Login struct {
	UserName string
	Password string
}

func NewLogin(userName string, password string) Login {
	return Login{
		UserName: userName,
		Password: password,
	}
}

func (l Login) Login(db *gorm.DB) (*Session, error) {

	// 1. Check Password
	user, err := GetUserByUserName(db, l.UserName)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(l.Password))
	if err != nil {
		return nil, err
	}
	// 2. Session Update/Create
	session, _ := GetSessionByUserULID(db, user.ID)

	accessToken, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	expired := time.Now().Add(time.Hour * 6)

	if session.ID == "" {
		session := NewSession(user.ID, string(accessToken), expired)
		err = session.Create(db)
		if err != nil {
			return nil, err
		}

		return &session, nil

	} else {
		session.AccessToken = string(accessToken)
		session.Expired = expired
		_, err := UpdateSession(db, *session)
		if err != nil {
			return nil, err
		}

		return session, nil
	}
}
