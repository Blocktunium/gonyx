package http

import "github.com/gofiber/fiber/v2"

// Mark: Definitions

// Server struct
type Server struct {
	config ServerConfig
	app    *fiber.App
}

// init - Server Constructor - It initializes the server
func (s *Server) init(config ServerConfig) error {
	s.config = config
	s.app = fiber.New(fiber.Config{
		Prefork: false,
	})
	return nil
}

// MARK: Public functions

// NewServer - create a new instance of Server and return it
func NewServer(config ServerConfig) (*Server, error) {
	server := &Server{}
	err := server.init(config)
	if err != nil {
		return nil, NewCreateServerErr(err)
	}
	return server, nil
}

// Start - start the server and listen to provided address
func (s *Server) Start() error {
	err := s.app.Listen(s.config.ListenAddress)
	if err != nil {
		return NewStartServerErr(s.config.ListenAddress, err)
	}
	return nil
}

// Stop - stop the server
func (s *Server) Stop() error {
	err := s.app.Shutdown()
	if err != nil {
		return NewShutdownServerErr(err)
	}
	return nil
}
