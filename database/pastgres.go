package database

import (
	"../config"
	"context"
	"fmt"
	"github.com/go-pg/pg"
)

var psg *Postgres
var isTx bool = false

type Postgres struct {
	Con *pg.DB
	Tx  *pg.Tx
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func GetConnection() *Postgres {
	return psg
}

func BeginTx() (*Postgres, error) {
	p := GetConnection()
	tx, err := p.Con.Begin()
	if err != nil {
		return p, err
	}

	// Rollback tx on error.
	defer tx.Rollback()

	isTx = true

	p.Tx = tx

	return p, err
}

func CommitTx() {
	p := GetConnection()
	p.Tx.Commit()
	isTx = false
}

func IsTx() bool {
	return isTx
}

func Initialize(env string) *Postgres {
	cf := config.LoadConfiguration(env)

	dbHost := cf.Database.DbHost
	dbUser := cf.Database.DbUser
	dbName := cf.Database.DbName
	dbSSLMode := cf.Database.DbSsl
	dbPass := cf.Database.DbPassword

	op, _ := pg.ParseURL("postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":5432/" + dbName + "?sslmode=" + dbSSLMode)
	db := pg.Connect(op)
	db.AddQueryHook(dbLogger{}) // logger db query

	//log.Println(db)

	psg = &Postgres{Con: db}

	return psg
}
