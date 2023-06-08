package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// CreateBool creates a pointer to a bool with a given value for one-line usage
func CreateBool(value bool) *bool {
	return &value
}

// WaitForShutdown waits for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
func WaitForShutdown(httpServer *http.Server, isShuttingDown *bool) {
	quit := make(chan os.Signal, 1)

	// syscall.SIGINT  --> kill -2 --> CTRL-C
	// syscall.SIGTERM --> kill -9 --> k8s
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// and let's go!
	log.Infof("shutdown: started")
	*isShuttingDown = true

	// shutdown http-server
	if httpServer != nil {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Fatalf("gin: %s", err)
		}
		log.Infof("gin: shutdown completed")
	}

	log.Infof("shutdown: completed")
}
