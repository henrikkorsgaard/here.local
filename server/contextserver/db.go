package contextserver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

//for now we will not use a database.
func initSqliteDB() {
	var err error
	db, err = sql.Open("sqlite3", "db/here.local.sqlite.db")
	if err != nil {
		fmt.Println("err!")
		fmt.Println(err)
	}

	stmt := `CREATE TABLE IF NOT EXISTS readings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		locationmac TEXT NOT NULL,
		devicehash TEXT NOT NULL,
		signal INT,
		timestamp INT
	)`
	_, err = db.Exec(stmt)

	if err != nil {
		fmt.Println(err)
	}

}

func insertReading(locationmac string, devicehash string, signal int, timestamp time.Time) {
	stmt, err := db.Prepare("INSERT INTO readings(locationmac, devicehash, signal, timestamp) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(locationmac, devicehash, signal, timestamp.Unix())

	if err != nil {
		fmt.Println("Error inserting reading into database!")
		fmt.Println(err)
	}
}
