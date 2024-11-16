package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knakul853/accessmesh/internal/api"
	"github.com/knakul853/accessmesh/internal/config"
	"github.com/knakul853/accessmesh/internal/store"
	"github.com/knakul853/accessmesh/pkg/enforcer"
)

func main() {
	cfg := config.Load()

	db, err := store.NewMongoStore(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	enforcer, err := enforcer.NewCasbinEnforcer(db)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	api.SetupRoutes(router, db, enforcer)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
