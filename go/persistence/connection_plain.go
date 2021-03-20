// +build !dqlite

package persistence

import (
	"database/sql"
	"path"

	log "github.com/sirupsen/logrus"
)

// GetConnection as
func GetConnection(dataDir string, sync string, certificate string, key string, join []string) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", path.Join(dataDir, "k8l.db"))
	if err != nil {
		log.Fatal("FATAL: opening db ", err)
	}

	cleanup := func() {
		db.Close()
	}
	return db, cleanup
}
