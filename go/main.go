package main

import (
	"flag"
	"net/http"
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
	var staticDir string
	var certificate string
	var key string
	var verbose bool
	listen := flag.String("listen", ":9090", "address used to expose main API")
	flag.StringVar(&dataDir, "data", "/tmp/fuck", "data dir")
	flag.StringVar(&staticDir, "static", "./static", "Static data dir dir")
	flag.StringVar(&certificate, "cert", "cluster.crt", "tls certificate file")
	flag.StringVar(&key, "key", "cluster.key", "tls  key file")
	flag.StringVar(&sync, "sync", "", "listening address for internal database replication default if empty is first interface :9000")
	seed := flag.String("seed", "", "database addresses of existing nodes")
	flag.BoolVar(&verbose, "verbose", false, "verbose log")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	var join []string = []string{}
	if *seed != "" {
		join = []string{*seed}
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal(err)
	}

	conn, cleanupFunc := p.GetConnection(dataDir, sync, certificate, key, join)
	defer cleanupFunc()
	repository := p.LogRepository{Connection: conn}
	repository.Setup()
	p.STORAGE.LogRepository = &repository
	gin.DisableConsoleColor()

	r := api.NewRouter()

	//r.Use(injectRepository(&repository))
	r.Use(gin.Recovery())

	r.Use(static.Serve("/static", static.LocalFile(staticDir, false)))
	r.LoadHTMLGlob(staticDir + "/templates/*.tmpl.html")

	r.POST("/_bulk", api.BulkHandler)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"title": "Home",
		})
	})
	log.Info("listen to: ", *listen)
	r.Run(*listen)

}
