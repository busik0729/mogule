package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
		CREATE TABLE mail(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			text text NOT NULL,
			template_id UUID,
			client_ids  text[] DEFAULT array[]::text[],
			bid SERIAL,
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_mail PRIMARY KEY ( id )
		);
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
			DROP TABLE mail;
		`)
		return err
	})
}
