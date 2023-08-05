package migration

import (
	"log"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Migration struct {
	Process func(session *r.Session) error
	Version int
}

func NewMigration(process func(session *r.Session) error, version int) Migration {
	return Migration{
		Process: process,
		Version: version,
	}
}

func NewMigrations() []Migration {
	migrations := []Migration{}

	migrations = append(migrations, NewMigration(migrateV1, 0)) // migration V1: Add user_v1 and session_v1 tables.
	return migrations
}

func Migrate(session *r.Session) {
	migrations := NewMigrations()

	for _, mig := range migrations {
		if mig.Version < version {
			err := mig.Process(session)
			if err != nil {
				log.Println(err)
				panic("Migration was stopped because the error was occurred.")
			}
		}
	}

}
