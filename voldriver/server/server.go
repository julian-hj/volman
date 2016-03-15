package server

import (
	"os"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"
)

type Server struct {
}

func New(logger lager.Logger, backend voldriver.Backend, transport, address string) *Server {
	return &Server{}
}

func (s *Server) Run(signal <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)
	return nil
}
