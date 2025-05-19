package main

import (
	"fmt"
	"log"
)

// ``type LoadBalancer struct {
// 	iter iterator.Iterator
// }

func main() {
	lb_config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}
	fmt.Println(lb_config)
	// router := mux.NewRouter()

	//1.  health checking for all existing services.

}
