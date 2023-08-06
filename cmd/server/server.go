package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/transparentt/login-server/config"
	"github.com/transparentt/login-server/pkg/rethinkdb/logic"
	"github.com/transparentt/login-server/pkg/rethinkdb/migration"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type SignUp struct {
	UserName string `json:"user_name"`
	PassWord string `json:"password"`
}

func main() {
	config := config.LoadConfig()

	session, err := r.Connect(r.ConnectOpts{
		Address:  config.Address,
		Database: config.Database,
	})
	if err != nil {
		log.Fatalln(err)
	}

	migrate(session)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", health)
	e.POST("/users", func(c echo.Context) error {
		return signUp(c, session)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8000"))

}

// Handler
func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func signUp(c echo.Context, rSession *r.Session) error {
	signUp := new(SignUp)
	if err := c.Bind(signUp); err != nil {
		return err
	}

	existing, err := logic.GetUserByUserName(rSession, signUp.UserName)
	if err != nil {
		return err
	}

	if existing.ID != "" {
		return c.String(http.StatusNotAcceptable, "NG")
	}

	user, err := logic.NewUser(signUp.UserName, signUp.PassWord)
	if err != nil {
		return err
	}

	_, err = user.Create(rSession)
	if err != nil {
		return err
	}

	return c.String(http.StatusCreated, "OK")
}

func migrate(session *r.Session) {
	migration.Migrate(session)
}
