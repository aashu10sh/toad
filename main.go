package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ashshelby/toad/entities"
)

func main() {
	servers := []entities.Server{
		entities.NewSimpleServer("http://api.1"),
		entities.NewSimpleServer("http://api.2"),
		entities.NewSimpleServer("http://api.3"),
	}

	loadBalancer, error := entities.InitializeLoadBalancer(":8080", servers)
	if error != nil {
		os.Exit(1)
	}

	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		loadBalancer.ServeProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Load Balancer  Running on %s\n", loadBalancer.Port)
	http.ListenAndServe(loadBalancer.Port, nil)
}
