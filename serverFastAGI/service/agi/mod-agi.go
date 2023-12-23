package agi

import (
	"net"
	"sync"

	"github.com/pkg/errors"
)

// AGIServer representa un servidor AGI
type AGIServer struct {
	listener net.Listener
	handler  HandlerFunc
	stopCh   chan struct{}
	wg       sync.WaitGroup
	doneCh   chan struct{} // nuevo canal para señalizar la finalización de las goroutines
}

// NewAGIServer crea un nuevo servidor AGI
func NewAGIServer(addr string, handler HandlerFunc) (*AGIServer, error) {
	if addr == "" {
		addr = "localhost:4573"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to server") //no se pudo vincular al servidor
	}

	return &AGIServer{
		listener: l,
		handler:  handler,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}), // inicializar el nuevo canal
	}, nil
}

// Listen inicia el servidor AGI
func (s *AGIServer) Listen() {
	defer close(s.doneCh) // cerrar el canal de finalización cuando todas las goroutines hayan terminado
listenerLoop:
	for {
		select {
		case <-s.stopCh:
			break listenerLoop
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			// Aquí puedes manejar el error o simplemente salir
			break listenerLoop
		}

		// Manejar la conexión en una goroutine separada
		s.wg.Add(1)
		go func(conn net.Conn) {
			defer s.wg.Done()
			defer conn.Close()
			s.handler(NewConn(conn))
		}(conn)
	}
}

// Cerrar cierra el servidor AGI
func (s *AGIServer) Close() {
	close(s.stopCh)
	s.listener.Close()
	<-s.doneCh // esperar a que todas las goroutines hayan terminado
}

// Modificado Listen para utilizar AGIServer
func ListenMod(addr string, handler HandlerFunc) (*AGIServer, error) {
	server, err := NewAGIServer(addr, handler)
	if err != nil {
		return nil, err
	}

	// Iniciar el servidor en una goroutine
	go server.Listen()

	return server, nil
}
