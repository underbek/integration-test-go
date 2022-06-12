package server

import (
	"net/http"

	"github.com/AndreyAndreevich/articles/user_service/handler"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
}

func New(addr string, h *handler.Handler) *Server {
	r := mux.NewRouter()
	r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)

	server := &http.Server{
		Handler: r,
		Addr:    addr,
	}

	return &Server{
		server: server,
	}
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}
