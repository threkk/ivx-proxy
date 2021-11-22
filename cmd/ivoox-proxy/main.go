package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/threkk/ivoox-proxy/internal/app"
)

var port int
var baseAddress string

func init() {
	flag.IntVar(&port, "port", 3000, "Specify an alternate port [default: 3000]")
	flag.StringVar(&baseAddress, "base", "localhost", "Specify the base address. This is import to generate the right URL [default: localhost]")
}

func main() {
	ivp := app.NewApp()
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
