package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"docqube.de/bookkeeper/pkg/database"
	categoryHandler "docqube.de/bookkeeper/pkg/services/category/handler"
	transactionHandler "docqube.de/bookkeeper/pkg/services/transaction/handler"
	"docqube.de/bookkeeper/pkg/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	port           = os.Getenv("PORT")
	isShuttingDown = utils.NewBool(false)
)

func main() {
	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	// perform database migration to newest version
	log.Debug("starting database migration...")
	err := database.DBMigration("file://.sql/migrations")
	if err != nil {
		log.Fatalf("db migration failed: %s", err)
	}
	log.Debug("database migration finished")

	g := gin.New()

	// router groups
	v1 := g.Group("/api/v1")
	v1.Use(gzip.Gzip(gzip.DefaultCompression))

	// register handlers
	_ = transactionHandler.NewHandler(v1)
	_ = categoryHandler.NewHandler(v1)

	// Check if any port is set, otherwise use 8080
	if len(port) == 0 {
		port = "8080"
	}

	// Check connectivity to the database
	db, err := database.GetConnection()
	if err != nil {
		log.Fatalf("Can't get database: %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't ping database: %s", err)
	}
	g.GET("/healthz/:probe", func(c *gin.Context) {
		probe := c.Param("probe")
		err = db.Ping()
		if err != nil {
			log.Errorf("%s: db ping error: %s", probe, err)
			c.AbortWithStatus(503)
			return
		}
		c.Status(200)
	})

	// build and start http-server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: g,
	}
	go func() {
		log.Infof("Starting http-server...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	// wait for shutdown
	utils.WaitForShutdown(httpServer, isShuttingDown)
	log.Infof("Goodbye!")
}
