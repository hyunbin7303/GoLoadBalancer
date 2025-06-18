package serverpool

import (
	"HpLoadBalancer/lb/server"
	"context"
	"log"
	"sync"
	"time"
)

type ServerPool struct {
	Servers []*server.Server
	mux     sync.RWMutex
	curr    int
}

func (sp *ServerPool) Rotate() *server.Server {
	sp.mux.Lock()
	sp.curr = (sp.curr + 1) % len(sp.Servers)
	sp.mux.Unlock()
	return sp.Servers[sp.curr]
}

func (sp *ServerPool) GetNextServer() *server.Server {
	for i := 0; i < len(sp.Servers); i++ {
		nextPeer := sp.Rotate()
		if nextPeer.Alive {
			return nextPeer
		}
	}
	return nil
}

func (sp *ServerPool) AddServer(server *server.Server) {
	// sp.Rotate().mux.Lock()
	sp.Servers = append(sp.Servers, server)
	// sp.mux.Unlock()
}

func LauchHealthCheck(ctx context.Context, sp ServerPool) {
	t := time.NewTicker(time.Second * 10)
	log.Println("Starting Health check.")
	for {
		select {
		case <-t.C:
			go HealthCheck(ctx, sp)
		case <-ctx.Done():
			log.Println("Closing health check.")
			return
		}
	}
}
