package initialise

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//for now we will not use a database.
func initDB() {
	db, err := sql.Open("sqlite3", "db/locations.sqlite.db")
	if err != nil {
		fmt.Println("err!")
		fmt.Println(err)
	}

	stmt := "SELECT name FROM sqlite_master WHERE type='table' AND name='locations';"
	_, err := db.Exec(stmt)

	var name string
	err = row.Scan(&name)
	if err == sql.ErrNoRows {
		fmt.Println("no fecking table exists -- better make one kiddo!")
	}

}
