// +build dqlite

package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
)

// GetConnection as
func GetConnection(dataDir string, sync string, join []string) (*sql.DB, func()) {
	logFunc := func(l client.LogLevel, format string, a ...interface{}) {
		log.Printf(fmt.Sprintf("%s: %s\n", l.String(), format), a...)
	}
	app, err := app.New(dataDir, app.WithAddress(sync), app.WithCluster(join), app.WithLogFunc(logFunc))
	if err != nil {
		log.Fatal("Cannot create app", err)
	}

	if err := app.Ready(context.Background()); err != nil {
		log.Fatal("app is not ready", err)
	}

	conn, err := app.Open(context.Background(), "k8l") // TODO: replace with hash based on sync listening address
	if err != nil {
		log.Fatal("Cannot open  database", err)
	}
	cleanup := func() {
		conn.Close()
		app.Handover(context.Background())
		app.Close()
	}
	return conn, cleanup
}
