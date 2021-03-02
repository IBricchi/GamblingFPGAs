package main

import (
	"context"
	"flag"
	"log"

	"github.com/IBricchi/GamblingFPGAs/server"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	var httpPort = flag.String("httpPort", "3000", "Port for serving http server")
	flag.Parse()

	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("buna: failed to create zap logger: %v\n", err)
	}
	defer logger.Sync()

	r := chi.NewRouter()

	httpServer, err := server.OpenHttpServer(ctx, logger, r)
	if err != nil {
		logger.Fatal("server: failed to open http server", zap.Error(err))
	}
	defer httpServer.Close()
	logger.Info("server: opened http server")

	if err := httpServer.Serve(*httpPort); err != nil {
		logger.Fatal("server: failed to serve http server", zap.Error(err))
	}
}
