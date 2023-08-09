package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/transparentt/login-server/pkg/rethinkdb/logic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SignUp struct {
	UserName string `json:"user_name"`
	PassWord string `json:"password"`
}

type Login struct {
	UserName string `json:"user_name"`
	PassWord string `json:"password"`
}

func main() {

	dsn := "host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&logic.User{})
	db.AutoMigrate(&logic.Session{})

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", health)
	e.POST("/users", func(c echo.Context) error {
		return signUp(c, db)
	})
	e.POST("/login", func(c echo.Context) error {
		return login(c, db)
	})
	e.GET("/secret", func(c echo.Context) error {
		return secret(c, db)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8000"))

}

// Handler
func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func signUp(c echo.Context, db *gorm.DB) error {
	signUp := new(SignUp)
	if err := c.Bind(signUp); err != nil {
		return err
	}

	existing, _ := logic.GetUserByUserName(db, signUp.UserName)

	if existing.ID != "" {
		return c.String(http.StatusNotAcceptable, "NG")
	}

	user, err := logic.NewUser(signUp.UserName, signUp.PassWord)
	if err != nil {
		return err
	}

	err = user.Create(db)
	if err != nil {
		return err
	}

	return c.String(http.StatusCreated, "OK")
}

func login(c echo.Context, db *gorm.DB) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		return err
	}

	newLogin := logic.NewLogin(login.UserName, login.PassWord)
	session, err := newLogin.Login(db)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = session.AccessToken
	cookie.Expires = session.Expired
	c.SetCookie(cookie)

	cookie2 := new(http.Cookie)
	cookie2.Name = "user_ul_id"
	cookie2.Value = session.UserULID
	cookie2.Expires = session.Expired
	c.SetCookie(cookie2)

	return c.String(http.StatusOK, "OK")
}

func secret(c echo.Context, db *gorm.DB) error {
	user_ulid, err := c.Cookie("user_ul_id")
	if err != nil {
		return err
	}

	access_token, err := c.Cookie("access_token")
	if err != nil {
		return err
	}

	session, err := logic.CheckSession(db, user_ulid.Value, access_token.Value)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = session.AccessToken
	cookie.Expires = session.Expired
	c.SetCookie(cookie)

	cookie2 := new(http.Cookie)
	cookie2.Name = "user_ul_id"
	cookie2.Value = session.UserULID
	cookie2.Expires = session.Expired
	c.SetCookie(cookie2)

	return c.String(http.StatusOK, "Secret OK")
}
