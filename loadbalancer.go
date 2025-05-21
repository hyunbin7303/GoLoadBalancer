package main

// type ServerPool interface {
// 	GetServers() []Server
// 	GetNextValidPeer() Server
// 	AddServer(Server)
// 	RemoveServer(Server)
// 	GetServerPoolSize() int
// }

// type LoadBalancer interface {
// 	Serve(http.ResponseWriter, *http.Request)
// }

// type LoadBalancer struct {
// 	Port     int
// 	LastPort int
// 	Count    int
// 	Servers  []*Server
// 	Config   Config
// }

// func (lb *LoadBalancer) getNextServer() *Server {
// 	lb.Count++
// 	server := lb.Servers[lb.Count % len(lb.Servers)]
// 	if !server.IsAlive(lb.Config.
// }
