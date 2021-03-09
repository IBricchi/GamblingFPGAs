package server

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
)

func (h *HttpServer) routes(ctx context.Context) error {
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

		// Post
		r.Post("/test/dynamic", h.handlePostDynamicTest(ctx))
		r.Post("/poker/openGame", h.handlePokerOpenGame())
		r.Patch("/poker/joinGame", h.handlePokerJoinGame())
	})

	// Not found
	h.router.NotFound(h.handleNotFound())

	return nil
}
