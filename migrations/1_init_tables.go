package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`

		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE users(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			username character varying(255) NOT NULL,
			name character varying(255) NOT NULL,
			surname character varying(255) NOT NULL,
			password character varying(255) NOT NULL,
			role integer NOT NULL,
			avatar character varying(255) DEFAULT 'assets/images/avatars/profile.jpg',
			created_at	 timestamp with time zone,
			updated_at	 timestamp with time zone,
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_users PRIMARY KEY ( id )
		);

		CREATE TABLE device(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			user_id UUID NOT NULL,
			platform	 character varying(255),
			device_uuid	 character varying(255),
			model	 character varying(255),
			serial	 character varying(255),
			version_os character varying(255) NOT NULL,
			version_app	 character varying(255) NOT NULL,
			access_token	 text NOT NULL UNIQUE,
			refresh_token	 text NOT NULL UNIQUE,
			expired_in	 bigint NOT NULL,
			created_at	 timestamp with time zone,
			updated_at	 timestamp with time zone,
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_device PRIMARY KEY ( id )
		);

		INSERT INTO "users" (username, name, surname, password, role) VALUES ('admin', 'admin', 'admin', '$2a$14$KE5Vax1ui6F54Sq43sgLsedw3HD2K4OMgpFvzi/gW5kXiU2uoCNii', 1);
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
			DROP TABLE users;
			DROP TABLE device;
		`)
		return err
	})
}
