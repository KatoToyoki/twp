package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jykuo-love-shiritori/twp/db"
	"github.com/jykuo-love-shiritori/twp/pkg/constants"
	"github.com/jykuo-love-shiritori/twp/pkg/router"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	logger, err := zap.NewProduction()
	if err != nil {
		e.Logger.Fatal(err)
	}
	RegisterLogger(e, logger)

	db, err := db.NewDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	RegisterFrontend(e)

	router.RegisterApi(e, db, logger.Sugar())

	if os.Getenv("TWP_ENV") == constants.DEV.String() {
		router.RegisterDocs(e)
	}
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
