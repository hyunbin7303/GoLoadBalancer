package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second) // Simulate processing time
	w.Write([]byte("Hello, world!"))
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer timeTrack(startTime, "Request to "+r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func main() {
	args := os.Args
	servername := args[1]
	portnum := args[2]
	homeHandler := func(w http.ResponseWriter, req *http.Request) {
		l := log.New(os.Stdout, "[Server 1] ", log.Ldate|log.Ltime)
		l.Printf(" Running .... ")
		fmt.Println("Testing home handler")
		io.WriteString(w, "Hi there. \n")
	}

	http.Handle("/", loggingMiddleware(http.HandlerFunc(homeHandler)))
	http.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler)))
	fmt.Printf("Starting %s .... at port %s\n", servername, portnum)
	log.Fatal(http.ListenAndServe(portnum, nil))
}
