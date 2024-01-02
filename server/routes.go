package server

import "go-dk/handlers"

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)

	handlers.FrontPage(s.mux)
}
