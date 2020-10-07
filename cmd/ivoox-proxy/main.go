package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/threkk/ivoox-proxy/internal/app"
)

func main() {
	ivp := app.NewApp("localhost:3000")
	srv := &http.Server{
		Handler:        ivp,
		Addr:           ":3000",
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	go func() {
		srv.Shutdown(ctx)
	}()

	<-ctx.Done()

	os.Exit(1)
}
