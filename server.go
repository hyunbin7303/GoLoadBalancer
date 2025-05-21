package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// type Server interface {
// 	SetAlive(bool)
// 	IsAlive() bool
// 	GetURL() *url.URL
// 	GetActiveConnections() int
// 	Serve(http.ResponseWriter, *http.Request)
// }

// type ServerPool interface {
// 	GetServers() []Server
// 	GetNextValidPeer() Server
// 	AddServer(Server)
// 	RemoveServer(Server)
// 	GetServerPoolSize() int
// }

type ServerPool struct {
	servers []Server
}

type Server struct {
	healthCheckPath string
	Address         *url.URL
	// IsAlive         bool
	ServerMux    http.Handler // or use sync.RwMutex
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) SetAlive(alive bool) {
	// TODO
}

func (s *Server) IsAlive() bool {
	// TODO
	return false
}

func (s *Server) Serve(w http.ResponseWriter, req *http.Request) {
	s.ReverseProxy.ServeHTTP(w, req)
}

func NewDevServer(addr string) *Server {
	serverUrl, _ := url.Parse("http://" + addr)
	// utils.OnErrorPanic(err, "Invalid server addr")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Response from server: %v\n", addr)
	})

	return &Server{
		healthCheckPath: serverUrl.String(),
		Address:         serverUrl,
		ServerMux:       mux,
		ReverseProxy:    NewSingleHostReverseProxy(serverUrl),
	}

}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}
