package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"mogui.it/k8l/go/api"
	p "mogui.it/k8l/go/persistence"
)

func injectRepository(repository *p.LogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("repository", repository)
		c.Next()
	}
}

func main() {

	// Setup logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	var sync string
	var dataDir string
	var verbose bool
	listen := flag.String("listen", ":9090", "address used to expose main API")
	flag.StringVar(&dataDir, "data", "/tmp/fuck", "data dir")
	flag.StringVar(&sync, "sync", "127.0.0.1:9001", "listening address for internal database replication")
	seed := flag.String("seed", "", "database addresses of existing nodes")
	flag.BoolVar(&verbose, "verbose", false, "verbose log")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	var join []string = []string{}
	if *seed != "" {
		join = []string{*seed}
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal(err)
	}

	conn, cleanupFunc := p.GetConnection(dataDir, sync, join)
	defer cleanupFunc()
	repository := p.LogRepository{Connection: conn}
	repository.Setup()
	p.STORAGE.LogRepository = &repository
	gin.DisableConsoleColor()

	r := api.NewRouter()
	//r.Use(injectRepository(&repository))
	r.Use(gin.Recovery())

	r.Use(static.Serve("/", static.LocalFile("./static", false)))

	r.POST("/_bulk", api.BulkHandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	log.Info("listen to: ", *listen)
	r.Run(*listen)

}
