package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
			CREATE TABLE event(
				id UUID NOT NULL DEFAULT uuid_generate_v1(),
				event_name character varying(255) NOT NULL,
				data json NOT NULL,
				device_id UUID,
				type_event integer NOT NULL,
				status integer NOT NULL,
				bid SERIAL,
				created_at	 timestamp with time zone,
				updated_at	 timestamp with time zone,
				deleted_at	 timestamp with time zone,
	
				CONSTRAINT id_event PRIMARY KEY ( id )
			);

			ALTER TABLE device ADD COLUMN ws_id character varying(255);
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
			DROP TABLE event;
		`)
		return err
	})
}
