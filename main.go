package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"mogui.it/k8l/api"
	"mogui.it/k8l/persistence"
)

func injectDB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
	}
}

func main() {

	listen := flag.String("listen", ":9090", "where to listen to")
	dbpath := flag.String("dbpath", "./test.db", "Sqlite database path")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbpath)
	if err != nil {
		log.Fatal("FATAL: opening db ", err)
	}
	defer db.Close()
	persistence.Setup(db)

	fmt.Println("listen to: ", *listen)

	gin.DisableConsoleColor()

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(injectDB(db))
	r.POST("/_bulk", api.BulkHandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.Run(*listen)
}
