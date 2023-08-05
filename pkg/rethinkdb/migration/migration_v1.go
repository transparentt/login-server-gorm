package migration

import (
	"github.com/transparentt/login-server/config"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func migrateV1(session *r.Session) error {

	config := config.LoadConfig()

	_, err := r.DB(config.Database).TableCreate("user_v1").RunWrite(session)
	if err != nil {
		return err
	}

	_, err = r.DB(config.Database).TableCreate("session_v1").RunWrite(session)
	if err != nil {
		return err
	}

	return nil
}
