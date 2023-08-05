package migration

import (
	"log"
	"time"

	"github.com/transparentt/login-server/config"
	"github.com/transparentt/login-server/pkg/rethinkdb/logic"
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

	migrations = append(migrations, NewMigration(migrateV1, 0)) // migration V1: Add user_v1 table.
	migrations = append(migrations, NewMigration(migrateV2, 1)) // migration V2: Add session_v1 table.

	return migrations
}

func Migrate(session *r.Session) {

	if wantedVersion <= 1 {
		config := config.LoadConfig()

		_, err := r.DB(config.Database).TableCreate(logic.MigrationTable).RunWrite(session)
		if err != nil {
			panic("Cannot create migration_v1 table!")
		}

		migrationStatus := NewMigrationStatus(0)
		migrationStatus.Create(session)
	}

	currentStatus, err := GetMigrationStatus(session)
	if err != nil {
		panic("Cannot get the current status of migration!")
	}

	migrations := NewMigrations()
	for version := currentStatus.Version; version < wantedVersion; version++ {

		err := migrations[version].Process(session)
		if err != nil {
			log.Println(err)
			panic("Migration was stopped because the error was occurred.")
		}

	}

	currentStatus.Version = wantedVersion
	currentStatus.MigratedAt = time.Now()

	_, err = UpdateMigrationStatus(session, *currentStatus)
	if err != nil {
		panic("Cannot update the MigrationStatus.")
	}

}

type MigrationStatus struct {
	ID         string    `json:"id" rethinkdb:"id"`
	Version    int       `json:"version" rethinkdb:"version"`
	MigratedAt time.Time `json:"migrated_at" rethinkdb:"migrated_at"`
}

func NewMigrationStatus(version int) MigrationStatus {
	return MigrationStatus{
		Version:    version,
		MigratedAt: time.Now(),
	}
}

func (ms *MigrationStatus) Create(session *r.Session) (*MigrationStatus, error) {
	config := config.LoadConfig()

	ms.ID = logic.NewULID().String()

	_, err := r.DB(config.Database).Table(logic.MigrationTable).Insert(ms).RunWrite(session)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func UpdateMigrationStatus(session *r.Session, migrationStatus MigrationStatus) (*MigrationStatus, error) {
	config := config.LoadConfig()

	_, err := r.DB(config.Database).Table(logic.MigrationTable).Get(migrationStatus.ID).Update(migrationStatus).RunWrite(session)
	if err != nil {
		return nil, err
	}

	return &migrationStatus, nil
}

func GetMigrationStatus(session *r.Session) (*MigrationStatus, error) {
	config := config.LoadConfig()

	cursor, err := r.DB(config.Database).Table(logic.MigrationTable).Run(session)
	if err != nil {
		return nil, err
	}

	current := MigrationStatus{}
	cursor.One(&current)
	cursor.Close()

	return &current, nil
}
