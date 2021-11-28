package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/okmaki/node/log/internal/adapter"
	"github.com/okmaki/node/log/internal/core"
	"github.com/okmaki/node/log/internal/handle"
)

func section(text string) {
	core.Info("main", text)
}

func main() {

	core.Configure()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	host := "0.0.0.0"
	port := 3000

	// ------------------------------
	section("configuring db connection...")
	// ------------------------------

	dbHosts := []string{"127.0.0.1"}
	db, err := adapter.NewStorageAdapter(dbHosts)
	if err != nil {
		core.Fatal("main", "failed to configure db connection - %v", err)
	}
	defer db.Close()

	// ------------------------------
	section("configuring hub...")
	// ------------------------------

	onLog := func(log core.Log) {
	}

	hub := adapter.NewHubAdapter(onLog)
	defer hub.Shutdown()

	go func() {
		if err := hub.Start(); err != nil {
			core.Fatal("main", "hub unexpectedly stopped running - %v", err)
		}
	}()

	// ------------------------------
	section("configuring middleware and routing...")
	// ------------------------------

	router := mux.NewRouter()
	router.HandleFunc("/log", handle.Record(db, hub)).Methods(http.MethodPost)
	router.HandleFunc("/search", handle.Search(db)).Methods(http.MethodPost)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	router.HandleFunc("/inspect", handle.Inspect(upgrader, hub))

	// ------------------------------
	section("configuring server...")
	// ------------------------------

	addr := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// ------------------------------
	section("starting server...")
	// ------------------------------

	go func() {
		core.Info("main", "listening on %s", addr)
		if err := server.ListenAndServe(); err != nil {
			core.Error("main", "server unexpectedly stopped running - %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	// ------------------------------
	section("shutting down...")
	// ------------------------------

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	os.Exit(0)
}
