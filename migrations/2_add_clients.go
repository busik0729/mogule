package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
		CREATE TABLE client(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			surname character varying(255),
			thirdname character varying(255),
			email character varying(255),
			phone character varying(255),
			vk character varying(255),
			instagram character varying(255),
			facebook character varying(255),
			whatsapp character varying(255),
			telegram character varying(255),
			twitter character varying(255),
			odnoklassniki character varying(255),
			category_id UUID NOT NULL,
			last_com timestamp with time zone,
			result_com text,
			manager UUID NOT NULL,
			created_at	 timestamp with time zone,
			updated_at	 timestamp with time zone,
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_client PRIMARY KEY ( id )
		);

		CREATE TABLE category_client(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			title character varying(255) NOT NULL,

			CONSTRAINT id_category_client PRIMARY KEY ( id )
		);

		INSERT INTO "category_client" (title) VALUES ('Агенство'), ('Врач'), ('Юрист');
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
			DROP TABLE client;
			DROP TABLE category_client;
		`)
		return err
	})
}
