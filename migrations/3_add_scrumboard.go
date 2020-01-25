package main

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS list_seq START 1;
		CREATE SEQUENCE IF NOT EXISTS card_seq START 1;

		CREATE TABLE board(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			uri character varying(255) NOT NULL,
			settings json,
			pm UUID,
			bid SERIAL,
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_board PRIMARY KEY ( id )
		);

		CREATE TABLE list(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			board_id UUID NOT NULL,
			bid SERIAL,
			position numeric NOT NULL DEFAULT nextval('list_seq'::regclass),
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_list PRIMARY KEY ( id )
		);

		CREATE TABLE card(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			description character varying(255),
			id_attachment_cover character varying(255) DEFAULT '',
			subscribed boolean,
			checklists json,
			check_items integer,
			check_items_checked integer,
			due timestamp with time zone,
			list_id UUID NOT NULL,
			idMembers text[] DEFAULT array[]::text[],
			idLabels text[] DEFAULT array[]::text[],
			bid SERIAL,
			position numeric NOT NULL DEFAULT nextval('card_seq'::regclass),
			deleted_at	 timestamp with time zone,

			CONSTRAINT id_card PRIMARY KEY ( id )
		);

		CREATE TABLE label(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			color character varying(255),

			CONSTRAINT id_label PRIMARY KEY ( id )
		);

		CREATE TABLE attachment(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			src character varying(255) NOT NULL,
			time timestamp with time zone,
			type integer,
			card_id UUID NOT NULL,

			CONSTRAINT id_attachment PRIMARY KEY ( id )
		);

		CREATE TABLE comment(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			member_id UUID NOT NULL,
			card_id UUID NOT NULL,
			message text,
			time timestamp with time zone,
			bid SERIAL,

			CONSTRAINT id_comment PRIMARY KEY ( id )
		);

		CREATE TABLE activity(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			member_id UUID NOT NULL,
			message text,
			time timestamp with time zone,
			card_id UUID NOT NULL,

			CONSTRAINT id_activity PRIMARY KEY ( id )
		);

		CREATE TABLE checklist(
			id UUID NOT NULL DEFAULT uuid_generate_v1(),
			name character varying(255) NOT NULL,
			check_items_checked integer,
			check_items json,
			card_id UUID NOT NULL,

			CONSTRAINT id_checklist PRIMARY KEY ( id )
		);
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`
			DROP TABLE board;
			DROP TABLE list;
			DROP TABLE card;
			DROP TABLE member;
			DROP TABLE label;
			DROP TABLE attachment;
			DROP TABLE comment;
			DROP TABLE activity;
			DROP TABLE checklist;
		`)
		return err
	})
}
