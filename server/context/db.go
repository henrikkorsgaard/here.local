package context

import (
	"database/sql"
	"fmt"
	"time"
)

var db sql.DB

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
		vendor TEXT,
		signal INT,
		timestamp INT
	)`
	_, err = db.Exec(stmt)

	if err != nil {
		fmt.Println(err)
	}

}

type reading struct {
	id          int
	locationmac string
	devicehash  string
	vendor      string
	signal      int
	timestamp   time.Time
}

func (r *reading) Insert() {
	stmt, err := db.Prepare("INSERT INTO readings(locationmac, devicehash, vendor, signal, timestamp) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(r.locationmac, r.devicehash, r.vendor, r.signal, r.timestamp.Unix())

	if err != nil {
		fmt.Println("Error inserting reading into database!")
		fmt.Println(err)
	}
}
