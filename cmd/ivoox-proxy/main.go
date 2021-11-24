package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/threkk/ivoox-proxy/internal/app"
)

var port int
var baseURL string
var user string
var password string
var userFile string

func init() {
	flag.IntVar(&port, "port", 3000, "Specify an alternate port [default: 3000]")
	flag.StringVar(&baseURL, "base", "", "Specify a custom base address.")
	flag.StringVar(&user, "user", "", "Username to secure the application.")
	flag.StringVar(&password, "password", "", "Password to secure the application. Mandatory if \"user\" is defined. Ignored otherwise.")
}

func main() {
	// Parse the flags and validate them.
	flag.Parse()
	if user != "" && password == "" {
		log.Fatal("If the user is defined, it needs a password")
	}

	// Server init.
	ivp := app.NewApp(baseURL, user, password)
	srv := &http.Server{
		Handler:        ivp,
		Addr:           fmt.Sprintf(":%d", port),
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shutdown
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	go func() {
		srv.Shutdown(ctx)
	}()

	<-ctx.Done()

	os.Exit(1)
}
