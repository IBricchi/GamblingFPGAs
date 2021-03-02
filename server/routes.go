package server

func (h *HttpServer) routes() {
	// Get
	h.router.Get("/test", h.handleGetStaticTest())

	// Post
	h.router.Post("/test/dynamic", h.handlePostDynamicTest())
}
