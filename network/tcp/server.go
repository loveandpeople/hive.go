package tcp

import (
	"fmt"
	"net"

	"github.com/loveandpeople/hive.go/events"
	"github.com/loveandpeople/hive.go/network"
	"github.com/loveandpeople/hive.go/syncutils"
)

type TCPServer struct {
	socket      net.Listener
	socketMutex syncutils.RWMutex
	Events      tcpServerEvents
}

func (srv *TCPServer) GetSocket() net.Listener {
	srv.socketMutex.RLock()
	defer srv.socketMutex.RUnlock()
	return srv.socket
}

func (srv *TCPServer) Shutdown() {
	srv.socketMutex.Lock()
	defer srv.socketMutex.Unlock()
	if srv.socket != nil {
		socket := srv.socket
		srv.socket = nil

		socket.Close()
	}
}

func (srv *TCPServer) Listen(bindAddress string, port int) *TCPServer {
	socket, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindAddress, port))
	if err != nil {
		println(fmt.Sprintf("TCP error: %s", err.Error()))
		srv.Events.Error.Trigger(err)

		return srv
	} else {
		srv.socketMutex.Lock()
		srv.socket = socket
		srv.socketMutex.Unlock()
	}

	srv.Events.Start.Trigger()
	defer srv.Events.Shutdown.Trigger()

	for srv.GetSocket() != nil {
		if socket, err := srv.GetSocket().Accept(); err != nil {
			if srv.GetSocket() != nil {
				println(fmt.Sprintf("TCP error: %s", err.Error()))
				srv.Events.Error.Trigger(err)
			}
		} else {
			peer := network.NewManagedConnection(socket)

			go srv.Events.Connect.Trigger(peer)
		}
	}

	return srv
}

func NewServer() *TCPServer {
	return &TCPServer{
		Events: tcpServerEvents{
			Start:    events.NewEvent(events.CallbackCaller),
			Shutdown: events.NewEvent(events.CallbackCaller),
			Connect:  events.NewEvent(network.ManagedConnectionCaller),
			Error:    events.NewEvent(events.ErrorCaller),
		},
	}
}
