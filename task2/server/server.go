package main

import (
	"context"
    "encoding/base64"
    "encoding/json"
	"errors"
	"log"
    "math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"task2/decode"
)

func Decode(w http.ResponseWriter, req *http.Request) {
    reqBody := decode.DecodeRequest{}

    err := json.NewDecoder(req.Body).Decode(&reqBody)

    if err != nil {
        http.Error(w, "Couldn't decode request string:" + err.Error(),
        http.StatusBadRequest)
        return
    }

    sDec, _ := base64.StdEncoding.DecodeString(reqBody.InputString)

    resp := decode.DecodeResponse{}
    resp.OutputString = string(sDec)
    
    body, err := json.Marshal(resp)

    if err != nil {
        http.Error(w, "Couldn't encode response string:" + err.Error(),
         http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(body)
}

func HardOp(w http.ResponseWriter, req *http.Request) {
    time.Sleep(time.Duration(rand.IntN(11) + 10) * time.Second)
    
    if rand.IntN(2) == 0 {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        w.WriteHeader(http.StatusOK)
    }
}

func Version(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("v1.0.1"))
}

func main() {
    http.HandleFunc("/decode", Decode)
    http.HandleFunc("/hard-op", HardOp)
    http.HandleFunc("/version", Version)

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

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("HTTP shutdown error:", err)
    }

    log.Println("Server shutdown")
}