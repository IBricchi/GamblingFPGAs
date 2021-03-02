package server

func (h *HttpServer) routes() {
	h.router.Get("/test", h.handleStaticTest())
}
