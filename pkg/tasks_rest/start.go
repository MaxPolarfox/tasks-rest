package tasks_rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log"
)

// Serve starts the HTTP server, blocking until the server exits.
func (s *Service) Start() {
	// listen to shutdown from the listen thread, before exiting the main thread
	shutDownChan := make(chan bool, 2)

	// listen to the appropriate signals, and notify a channel
	stopChan := make(chan os.Signal, 10)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.options.Port),
		Handler: s.Router,
	}

	server.RegisterOnShutdown(func() {
		transport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			panic("Cannot cast http.DefaultTransport to *http.Transport")
		}
		transport.DisableKeepAlives = true
		transport.CloseIdleConnections()
		server.SetKeepAlivesEnabled(false)
		log.Println("RegisterOnShutdownCompleted")
	})

	go func() {
		log.Println("serving HTTP", "port", s.options.Port)

		// https://golang.org/pkg/net/http/#Server.Shutdown says when
		// server.Shutdown() is called, ErrServerClosed will be returned,
		// so we're only capturing other errors
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("serverError", "err", err)
		}

		log.Println("serverExited")
		shutDownChan <- true
	}()

	<-stopChan // wait for a signal to exit
	log.Println("shutting down HTTP server")

	// shutdown the server by gracefully draining connections
	// See https://golang.org/pkg/net/http/#Server.Shutdown for more details.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		shutDownChan <- true
		log.Fatalf("shutdownError", "err", err)
	}

	<-shutDownChan
	log.Println("shutdownComplete")
}
