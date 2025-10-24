package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/anhtr13/synth-socket/socket/internal/conf"
	"github.com/anhtr13/synth-socket/socket/internal/socket"
)

func init() {
	conf.InitConnection()
}

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	socketServer := socket.NewSocketServer()

	handler := http.NewServeMux()
	handler.HandleFunc("/ws", socketServer.HandleWebSocket)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.PORT),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool, 1)

	go gracefulShutdown(httpServer, done)
	go socketServer.HandleQueue_RoomIO(done)
	go socketServer.HandleQueue_FriendIO(done)
	go socketServer.HandleQueue_Notification(done)

	log.Println("Server is running on port", httpServer.Addr)
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("Http server error: %s", err))
	}

	<-done

	conf.CloseAllConnections()
	log.Println("Graceful shutdown complete.")
}
