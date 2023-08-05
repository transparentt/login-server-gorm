package logic

import (
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Session struct {
	ID          string
	UserULID    string
	AccessToken string
	Expired     time.Time
}

func NewSession(userULID string, accessToken string, expired time.Time) Session {
	return Session{
		UserULID:    userULID,
		AccessToken: accessToken,
		Expired:     expired,
	}
}

func (s *Session) Create(session *r.Session) (*Session, error) {
	s.ID = NewULID().String()

	_, err := r.Table(SessionTable).Insert(s).RunWrite(session)
	if err != nil {
		return nil, err
	}

	return s, err
}
