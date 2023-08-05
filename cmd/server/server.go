package main

import (
	"log"

	"github.com/transparentt/login-server/config"
	"github.com/transparentt/login-server/pkg/rethinkdb/migration"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

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

}

func migrate(session *r.Session) {
	migration.Migrate(session)
}
