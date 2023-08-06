package logic

import (
	"github.com/transparentt/login-server/config"
	"golang.org/x/crypto/bcrypt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type User struct {
	ID             string `json:"id" rethinkdb:"id"`
	UserName       string `json:"user_name" rethinkdb:"user_name"`
	HashedPassword string `json:"hashed_password" rethinkdb:"hashed_password"`
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

func (user *User) Create(session *r.Session) (*User, error) {
	config := config.LoadConfig()
	user.ID = NewULID().String()

	_, err := r.DB(config.Database).Table(UserTable).Insert(user).RunWrite(session)
	if err != nil {
		return nil, err
	}

	return user, err

}

func GetUserByUserName(session *r.Session, userName string) (*User, error) {
	config := config.LoadConfig()

	cursor, err := r.DB(config.Database).Table(UserTable).Filter(r.Row.Field("user_name").Eq(userName)).Run(session)
	if err != nil {
		return nil, err
	}

	user := User{}
	cursor.One(&user)
	cursor.Close()

	return &user, nil
}
