package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pratikdevnal/Learning-Golang/internal/config"
)

func main(){
	fmt.Println("Server Started")

	cfg := config.MustLoad()

	router:= http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	server := http.Server{
    Addr: cfg.HTTPServer.Addr,
    Handler: router,
	}

	// log.Printf("Starting server on %s...", cfg.HTTPServer.Addr)
	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))
	
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func ()  {
		err := server.ListenAndServe()
		if err != nil {
    	log.Fatalf("failed to start server: %v", err)
	}
	}()

	<-done

	slog.Info("shutting down the server")  

	ctx, cancel :=context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err:=server.Shutdown(ctx)

	if err!= nil{
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
	
}