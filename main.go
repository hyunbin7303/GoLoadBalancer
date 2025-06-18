package main

import (
	"HpLoadBalancer/lb/server"
	"HpLoadBalancer/lb/serverpool"
	"HpLoadBalancer/lb/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func loadBalancing(w http.ResponseWriter, r *http.Request) {
	log.Println("Arrived Request : " + r.RequestURI)
	peer := serverPool.GetNextServer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service is not available", http.StatusServiceUnavailable)
}

var serverPool serverpool.ServerPool

func main() {
	lb_config, err := utils.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}
	fmt.Println(lb_config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, s := range lb_config.ServersPath {
		serverUrl, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		rp := httputil.NewSingleHostReverseProxy(serverUrl)
		backend := server.NewServer(serverUrl, rp, lb_config.HealthCheckPath)
		rp.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Fatalf("[%s] %s\n", serverUrl.Host, e)
			log.Fatal("Error handling the request", e)
			retries := utils.GetRetryFromContext(request)
			if retries < lb_config.RetryLimit {
				log.Println("Retrying... ")
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), 1, retries+1)
					rp.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}
			log.Println("Set Alive false.")
			backend.SetAlive(false) // after retries, mark this server as down.
		}

		serverPool.AddServer(backend)
		log.Printf("Server is configured : %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", lb_config.Port),
		Handler: http.HandlerFunc(loadBalancing),
	}

	// go serverPool.HealthCheck()
	go serverpool.LauchHealthCheck(ctx, serverPool)
	// go func() {
	// 	<-ctx.Done()
	// 	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// 	if err := server.Shutdown(shutdownCtx); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// go serverPool.HealthCheck()
	// go serverpool.HealthCheck(ctx, serverPool)
	// go func() {
	// 	<-ctx.Done()
	// 	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// 	if err := server.Shutdown(shutdownCtx); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	log.Printf("Load Balancer - activated with port :%d\n", lb_config.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen and server error :[%s] \n", err)
	}
}

//https://github.com/kasvith/simplelb/blob/master/main.go#L73
//https://medium.com/@leonardo5621_66451/building-a-load-balancer-in-go-1c68131dc0ef
