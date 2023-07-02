package main

import (
	"fmt"
	"net/http"
	"os"

	"docqube.de/bookkeeper/pkg/config"
	"docqube.de/bookkeeper/pkg/database"
	categoryHandler "docqube.de/bookkeeper/pkg/services/category/handler"
	intervalHandler "docqube.de/bookkeeper/pkg/services/interval/handler"
	transactionHandler "docqube.de/bookkeeper/pkg/services/transaction/handler"
	"docqube.de/bookkeeper/pkg/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	exitCode = 0

	isShuttingDown = utils.NewBool(false)
)

func main() {
	defer func() {
		os.Exit(exitCode)
	}()

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	config, err := config.LoadConfig()
	if err != nil {
		exitCode = 1
		log.Errorf("loading config: %s", err)
		return
	}
	log.Debugf("loaded config: %+v", config)

	db, err := database.InitializeDatabase(config)
	if err != nil {
		exitCode = 1
		log.Errorf("initializing database: %s", err)
		return
	}
	defer db.Close()
	log.Info("database initialized")

	g := gin.New()

	// router groups
	v1 := g.Group("/api/v1")
	v1.Use(gzip.Gzip(gzip.DefaultCompression))

	// register handlers
	_ = transactionHandler.NewHandler(v1, db)
	_ = categoryHandler.NewHandler(v1, db)
	_ = intervalHandler.NewHandler(v1, db)

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
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: g,
	}
	go func() {
		log.Infof("http-server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("unable to start http server listen: %s", err)
			os.Exit(1)
			return
		}
	}()

	// wait for shutdown
	utils.WaitForShutdown(httpServer, isShuttingDown)
	log.Infof("goodbye!")
}
