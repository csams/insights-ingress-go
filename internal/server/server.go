package server

import (
	"net/http"
)

type Server struct {
	CompletedConfig
}

type preparedServer struct {
	*Server
    Main http.Server
    Monitor http.Server
}

func New(c CompletedConfig) (*Server, error) {
	return &Server{
		CompletedConfig: c,
	}, nil
}

func (s *Server) PrepareRun() preparedServer {
	return preparedServer{s}
}

func (s preparedServer) Run() error {
	s.Log.V(0).Info("Listening on", "address", s.Options.Address)
	return http.ListenAndServe(s.Options.Address, s.RootHandler)
}
