package main

import (
	"database/sql"
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"mogui.it/k8l/api"
	. "mogui.it/k8l/persistence"
)

func injectRepository(repository *LogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("repository", repository)
	}
}

func main() {

	// Setup logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	listen := flag.String("listen", ":9090", "where to listen to")
	dbpath := flag.String("dbpath", "./test.db", "Sqlite database path")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbpath)
	if err != nil {
		log.Fatal("FATAL: opening db ", err)
	}
	defer db.Close()
	repository := LogRepository{Connection: db}
	repository.Setup()

	log.Info("listen to: ", *listen)

	gin.DisableConsoleColor()

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(injectRepository(&repository))

	r.POST("/_bulk", api.BulkHandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.Run(*listen)
}
