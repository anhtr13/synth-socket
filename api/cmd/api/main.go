package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anhtr13/synth-socket/api/internal/handler"
	"github.com/anhtr13/synth-socket/api/internal/conf"
)

func init() {
	conf.InitConnection()
}

func main() {
	err := loadDataToCache()
	if err != nil {
		log.Fatal("Cannot load data to cache:", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.PORT),
		Handler:      handler.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	forever := make(chan bool, 1)
	go gracefulShutdown(server, forever)

	log.Println("Server is running on port", server.Addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("Http server error: %s", err))
	}

	<-forever

	conf.CloseAllConnections()
	log.Println("Graceful shutdown complete.")
}
