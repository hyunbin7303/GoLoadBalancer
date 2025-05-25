package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	peer := serverPool.GetNextServer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service is not available", http.StatusServiceUnavailable)
}

var serverPool ServerPool

func main() {
	// lb_config, err := ReadConfig("config.yaml")
	// if err != nil {
	// 	log.Fatalf("read config error: %s", err)
	// }
	// fmt.Println(lb_config)

	serverConfigs := []struct {
		Address         string
		HealthCheckPath string
	}{
		{"http://localhost:5061", "/health"},
		{"http://localhost:5062", "/health"},
		{"http://localhost:5063", "/health"},
	}

	// pool := NewServerPool(serverConfigs)
	for _, server := range serverConfigs {
		u := server.Address
		serverUrl, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Fatal("[%s] %s\n", serverUrl.Host, e.Error())
			log.Fatal("Error handling the request", e)
			retries := GetRetryFromContext(request)
			if retries < 3 {

			}
		}

		serverPool.AddServer(&Server{
			Address:      serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Server is configured : %s\n", serverUrl)
	}

	// Creating http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", 8099),
		Handler: http.HandlerFunc(loadBalancing),
	}

	// go serverPool.LaunchHealthCheck(ctx, serverPool)

	log.Println("Load Balancer - activated with port :%d\n", 8099)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServer() Error", err)
	}

}

//https://github.com/leonardo5621/golang-load-balancer/blob/master/main.go
//https://github.com/kasvith/simplelb/blob/master/main.go#L73
//https://medium.com/@leonardo5621_66451/building-a-load-balancer-in-go-1c68131dc0ef
// https://github.com/kasvith/simplelb/blob/master/main.go
//https://github.com/Fuad28/load-balancer/blob/master/load_balancer.go
