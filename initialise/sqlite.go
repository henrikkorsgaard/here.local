package initialise

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//for now we will not use a database.
func initDB() sql.DB {
	db, err := sql.Open("sqlite3", "db/here.local.sqlite.db")
	if err != nil {
		fmt.Println("err!")
		fmt.Println(err)
	}

	stmt := `CREATE TABLE IF NOT EXISTS proximities (
		id TEXT PRIMARY KEY NOT NULL,
		location_id TEXT NOT NULL,
		device_id TEXT NOT NULL,
		signal INT,
		signals TEXT
	)`
	_, err = db.Exec(stmt)

	if err != nil {
		fmt.Println(err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS locations (
		mac TEXT PRIMARY KEY NOT NULL,
		ip TEXT NOT NULL,
		name TEXT NOT NULL,
		lastseen TEXT
	)`
	_, err = db.Exec(stmt)

	if err != nil {
		fmt.Println(err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT,
		vendor TEXT,
		public INT,
		lastseen TEXT
	)`
	_, err = db.Exec(stmt)

	if err != nil {
		fmt.Println(err)
	}

	return db
}
