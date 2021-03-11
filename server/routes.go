package server

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (h *HttpServer) routes(ctx context.Context) error {
	// General middleware
	h.router.Use(middleware.Logger)
	h.router.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"https://*", "http://*"}}))

	// Public routes
	h.router.Group(func(r chi.Router) {
		// Get
		r.Get("/public/test", h.handleGetStaticTest())

		// Post
		r.Post("/public/test/dynamic", h.handlePostDynamicTest(ctx))
	})

	// Private routes
	creds, err := h.db.getCreds(ctx)
	if err != nil {
		return fmt.Errorf("server: routes: failed to get creds: %w", err)
	}
	h.router.Group(func(r chi.Router) {
		r.Use(h.basicAuth("GamblingFPGAs-API", creds))

		// Get
		r.Get("/test", h.handleGetStaticTest())
		r.Get("/poker/openGameStatus", h.handlePokerGetGameOpenStatus())
		r.Get("/poker/activeGameStatus", h.handlePokerGetGameActiveStatus())

		// Post
		r.Post("/test/dynamic", h.handlePostDynamicTest(ctx))
		r.Post("/poker/openGame", h.handlePokerOpenGame())
		r.Post("/poker/joinGame", h.handlePokerJoinGame())
		r.Post("/poker/startGame", h.handlePokerStartGame())
		r.Post("/poker/terminateGame", h.handlePokerTerminateGame())
	})

	// Not found
	h.router.NotFound(h.handleNotFound())

	return nil
}
