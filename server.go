package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	healthCheckPath string
	Alive           bool
	Address         *url.URL
	mux             sync.RWMutex
	ReverseProxy    *httputil.ReverseProxy
}

func (s *Server) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

// func (s *Server) SetAlive(alive bool) {
// 	s.mux.Lock()
// 	s.Alive = alive
// 	s.mux.Unlock()
// }

// func singleJoiningSlash(a, b string) string {
// 	aslash := strings.HasSuffix(a, "/")
// 	bslash := strings.HasPrefix(b, "/")
// 	switch {
// 	case aslash && bslash:
// 		return a + b[1:]
// 	case !aslash && !bslash:
// 		return a + "/" + b
// 	}
// 	return a + b
// }

// func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
// 	targetQuery := target.RawQuery
// 	director := func(req *http.Request) {
// 		req.URL.Scheme = target.Scheme
// 		req.URL.Host = target.Host
// 		req.Host = target.Host
// 		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
// 		if targetQuery == "" || req.URL.RawQuery == "" {
// 			req.URL.RawQuery = targetQuery + req.URL.RawQuery
// 		} else {
// 			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
// 		}
// 	}
// 	return &httputil.ReverseProxy{Director: director}
// }
