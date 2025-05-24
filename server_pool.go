package main

import (
	"net/http"
	"sync"
)

func NewServerPool(serverConfigs []struct {
	Address         string
	HealthCheckPath string
}) *ServerPool {
	servers := make([]*Server, len(serverConfigs))
	for i, cfg := range serverConfigs {
		servers[i] = NewServer(cfg.Address, cfg.HealthCheckPath)
	}
	return &ServerPool{Servers: servers}
}

type ServerPool struct {
	Servers []*Server
	mu      sync.Mutex
	idx     int
}

func (sp *ServerPool) GetNextServer() *Server {
	sp.mu.Lock()
	server := sp.Servers[sp.idx%len(sp.Servers)]
	sp.idx++
	sp.mu.Unlock()
	return server
}

func (sp *ServerPool) AddServer(server *Server) {
	sp.mu.Lock()
	sp.Servers = append(sp.Servers, server)
	sp.mu.Unlock()
}

func (sp *ServerPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := sp.GetNextServer()
	server.ReverseProxy.ServeHTTP(w, r)
}
