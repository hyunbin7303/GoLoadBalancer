package main

import (
	"sync"
)

type ServerPool struct {
	Servers []*Server
	mux     sync.RWMutex
	curr    int
}

func (sp *ServerPool) Rotate() *Server {
	sp.mux.Lock()
	sp.curr = (sp.curr + 1) % len(sp.Servers)
	sp.mux.Unlock()
	return sp.Servers[sp.curr]
}

func (sp *ServerPool) GetNextServer() *Server {
	for i := 0; i < len(sp.Servers); i++ {
		nextPeer := sp.Rotate()
		if nextPeer.Alive {
			return nextPeer
		}
	}
	return nil
}

func (sp *ServerPool) AddServer(server *Server) {
	// sp.Rotate().mux.Lock()
	sp.Servers = append(sp.Servers, server)
	// sp.mux.Unlock()
}

// func HealthCheck(ctx context.Context, sp ServerPool) {
// 	// aliveChannel := make(chan bool, 1)
// 	// for _, server := range sp.Servers {
// 	// 	server := server
// 	// 	requestCtx, stop := context.WithTimeout(ctx, 10)
// 	// }
// }
