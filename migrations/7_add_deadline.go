package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
			ALTER TABLE board ADD COLUMN due timestamp with time zone;
			ALTER TABLE board ADD COLUMN client_id uuid;
			ALTER TABLE board ADD COLUMN daemon_status integer;
			ALTER TABLE card ADD COLUMN daemon_status integer;
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
		`)
		return err
	})
}
