package main

import (
	"log"
	"net/http"
)

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
		{"http://localhost:5001", "/health"},
		{"http://localhost:5002", "/health"},
		{"http://localhost:5003", "/health"},
	}

	pool := NewServerPool(serverConfigs)
	log.Println("Server pool started at :8099")
	if err := http.ListenAndServe(":8099", pool); err != nil {
		// pool.GetNextServer()
		log.Fatal(err)
	}

}
