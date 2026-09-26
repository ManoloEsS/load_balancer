package gateway

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GatewayServer struct {
	server *http.Server
	mux    *http.ServeMux
	proxy  *httputil.ReverseProxy
}

func NewGatewayServer(proxy *httputil.ReverseProxy, addr string) *GatewayServer {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &GatewayServer{
		server: srv,
		mux:    mux,
		proxy:  proxy,
	}

}

func (gw *GatewayServer) Setup() {
	gw.mux.Handle("/", gw.proxy)
}

func (gw *GatewayServer) Run() {
	gw.Setup()
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Println("HTTP gateway server listening on", gw.server.Addr)

		if err := gw.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}

		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("gateway server failed: %v", err)
		}
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gw.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if err := gw.server.Close(); err != nil {
			log.Printf("Forced shutdown failed: %v", err)
		}
	}

	log.Println("Server stopped")

}
