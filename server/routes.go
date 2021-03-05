package server

import "context"

func (h *HttpServer) routes(ctx context.Context) {
	// Get
	h.router.Get("/test", h.handleGetStaticTest())

	// Post
	h.router.Post("/test/dynamic", h.handlePostDynamicTest(ctx))

	// Not found
	h.router.NotFound(h.handleNotFound())
}
