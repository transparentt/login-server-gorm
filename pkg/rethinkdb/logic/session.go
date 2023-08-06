package logic

import (
	"time"

	"github.com/transparentt/login-server/config"
	"golang.org/x/crypto/bcrypt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Session struct {
	ID          string    `json:"id" rethinkdb:"id"`
	UserULID    string    `json:"user_ulid" rethinkdb:"user_ulid"`
	AccessToken string    `json:"access_token" rethinkdb:"access_token"`
	Expired     time.Time `json:"expired" rethinkdb:"expired"`
}

func NewSession(userULID string, accessToken string, expired time.Time) Session {
	return Session{
		UserULID:    userULID,
		AccessToken: accessToken,
		Expired:     expired,
	}
}

func (s *Session) Create(rSession *r.Session) (*Session, error) {
	config := config.LoadConfig()
	s.ID = NewULID().String()

	_, err := r.DB(config.Database).Table(SessionTable).Insert(s).RunWrite(rSession)
	if err != nil {
		return nil, err
	}

	return s, err
}

func GetSessionByUserULID(rSession *r.Session, userULID string) (*Session, error) {
	config := config.LoadConfig()

	cursor, err := r.DB(config.Database).Table(SessionTable).Filter(r.Row.Field("user_ulid").Eq(userULID)).Run(rSession)
	if err != nil {
		return nil, err
	}

	session := Session{}
	cursor.One(&session)
	cursor.Close()

	return &session, nil
}

func UpdateSession(rSession *r.Session, session Session) (*Session, error) {
	config := config.LoadConfig()

	_, err := r.DB(config.Database).Table(SessionTable).Get(session.ID).Update(session).RunWrite(rSession)
	if err != nil {
		return nil, err
	}

	return &session, nil
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

func (l Login) Login(rSession *r.Session) (*Session, error) {

	// 1. Check Password
	user, err := GetUserByUserName(rSession, l.UserName)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(l.Password))
	if err != nil {
		return nil, err
	}
	// 2. Session Update/Create
	session, err := GetSessionByUserULID(rSession, user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	expired := time.Now().Add(time.Hour * 6)

	if session == nil {
		session := NewSession(user.ID, string(accessToken), expired)
		_, err = session.Create(rSession)
		if err != nil {
			return nil, err
		}

		return &session, nil

	} else {
		session.AccessToken = string(accessToken)
		session.Expired = expired
		_, err := UpdateSession(rSession, *session)
		if err != nil {
			return nil, err
		}

		return session, nil
	}
}
