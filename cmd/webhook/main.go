package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kubevirt.io/irsa-mutation-webhook/internal/mutator"
	"kubevirt.io/irsa-mutation-webhook/pkg/config"
)

const (
	defaultTLSCertFile = "/etc/webhook/certs/tls.crt"
	defaultTLSKeyFile  = "/etc/webhook/certs/tls.key"
)

func main() {
	port := flag.Int("port", 8443, "Webhook server port")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	mutator, err := mutator.NewMutator(cfg)
	if err != nil {
		fmt.Printf("Failed to create mutator: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", mutator.HandleMutate)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	go func() {
		fmt.Printf("Starting webhook server on port %d\n", *port)
		if err := server.ListenAndServeTLS(defaultTLSCertFile, defaultTLSKeyFile); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
		os.Exit(1)
	}
}
