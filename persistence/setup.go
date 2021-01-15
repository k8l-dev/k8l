package persistence

import (
	"database/sql"
	"log"
)

// Setup the database, Create table and stuff
func Setup(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (	
		id INTEGER PRIMARY KEY AUTOINCREMENT,
			namespace_name TEXT NOT NULL,
			container_name TEXT NOT NULL,
			pod_name TEXT NOT NULL,
			container_image TEXT NOT NULL,
			timestamp TEXT NOT NULL,
			message TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal(err)
	}
}
