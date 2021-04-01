package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var db *sql.DB

type client struct {
	id   int
	name string
	pass string
}

type category struct {
	id          int
	description string
}

type budget struct {
	client_id int
	cat_id    int
	amount    sql.NullFloat64
}

type transaction struct {
	ID          int
	Client_id   int
	Cat_id      sql.NullInt64
	Amount      float64
	Balance     float64
	Description sql.NullString
	Time        time.Time
}

type stock struct {
	client_id int
	symbol    string
	quantity  int
}

// initializes the database. panics if a failure.
func init() {
	var err error
	// do NOT use in production
	connStr := "user=postgres password=? dbname=ecocosts sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}