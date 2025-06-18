package serverpool

import (
	"context"
	"log"
	"net"
	"net/url"
	"time"
)

func HealthCheck(ctx context.Context, sp ServerPool) {
	aliveChannel := make(chan bool, 1)
	for _, server := range sp.Servers {
		requestCtx, stop := context.WithTimeout(ctx, 3*time.Second)
		defer stop()
		status := "up"
		go isServerAlive(requestCtx, aliveChannel, server.Address)

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
		log.Printf("URL Status : %s for url : %s", status, server.Address)
	}
}
func isServerAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}
