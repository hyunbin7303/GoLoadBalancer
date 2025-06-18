package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	Attempts int = iota
	Retry
)

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

func loadBalancing(w http.ResponseWriter, r *http.Request) {
	log.Println("Arrived Request : " + r.RequestURI)
	peer := serverPool.GetNextServer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service is not available", http.StatusServiceUnavailable)
}

func IsServerAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		log.Fatal("Checking for IsServerAlive")
		log.Fatal(err)
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}

func HealthCheck(ctx context.Context, sp ServerPool) {
	aliveChannel := make(chan bool, 1)
	for _, server := range sp.Servers {
		server := server
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		go IsServerAlive(requestCtx, aliveChannel, server.Address)

		select {
		case <-ctx.Done():
			log.Println("Shutting down health check")
			return
		case alive := <-aliveChannel:
			server.Alive = alive
			if !alive {
				status = "down"
			}
		}
		log.Printf("URL Status : %s", status)
	}

}

var serverPool ServerPool

func main() {
	lb_config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}
	fmt.Println(lb_config)

	// pool := NewServerPool(serverConfigs)
	for _, server := range lb_config.ServersPath {
		serverUrl, err := url.Parse(server)
		if err != nil {
			log.Fatal(err)
		}
		rp := httputil.NewSingleHostReverseProxy(serverUrl)
		server := NewServer(serverUrl, rp)
		rp.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Fatalf("[%s] %s\n", serverUrl.Host, e)
			log.Fatal("Error handling the request", e)
			retries := GetRetryFromContext(request)
			if retries < lb_config.RetryLimit {
				log.Println("Retrying... ")
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					rp.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}
		}
		serverPool.AddServer(server)
		log.Printf("Server is configured : %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", 8099),
		Handler: http.HandlerFunc(loadBalancing),
	}

	go serverPool.HealthCheck()

	log.Printf("Load Balancer - activated with port :%d\n", 8099)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen and server error :[%s] \n", err)
	}
}

//https://github.com/leonardo5621/golang-load-balancer/blob/master/main.go
//https://github.com/kasvith/simplelb/blob/master/main.go#L73
//https://medium.com/@leonardo5621_66451/building-a-load-balancer-in-go-1c68131dc0ef
// https://github.com/kasvith/simplelb/blob/master/main.go
//https://github.com/Fuad28/load-balancer/blob/master/load_balancer.go
