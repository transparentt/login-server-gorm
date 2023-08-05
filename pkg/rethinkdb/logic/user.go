package logic

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type User struct {
	ID             string `json:"id" rethinkdb:"id"`
	UserName       string `json:"user_name" rethinkdb:"user_name"`
	HashedPassword string `json:"hashed_password" rethinkdb:"hashed_password"`
}

func NewUser(userName string, hashedPassword string) User {
	user := User{
		UserName:       userName,
		HashedPassword: hashedPassword,
	}
	return user
}

func (user *User) Create(session *r.Session) (*User, error) {
	user.ID = NewULID().String()
	_, err := r.Table(UserTable).Insert(user).RunWrite(session)
	if err != nil {
		return nil, err
	}

	return user, err

}
