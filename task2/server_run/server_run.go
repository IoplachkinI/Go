package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	srv "task2/server"
)

func main() {
    http.HandleFunc("/decode", srv.Decode)
    http.HandleFunc("/hard-op", srv.HardOp)
    http.HandleFunc("/version", srv.Version)

	const timeout = 10 * time.Second

    server := &http.Server{
        Addr: ":8080",
    }

    log.Println("Server running")

    go func() {
        if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            log.Fatal("HTTP server error:", err)
        }
        log.Println("Stopped serving new connections")
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("HTTP shutdown error:", err)
    }

    log.Println("Server shutdown")
}
