package entities

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type SimpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type LoadBalancer struct {
	Port            string
	roundRobinCount int
	servers         []Server
}

func InitializeLoadBalancer(port string, servers []Server) (*LoadBalancer, error) {
	return &LoadBalancer{
		roundRobinCount: 0,
		Port:            port,
		servers:         servers,
	}, nil
}

func NewSimpleServer(address string) *SimpleServer {
	url, error := url.Parse(address)
	if error != nil {
		os.Exit(1)
	}
	return &SimpleServer{address: address, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (server *SimpleServer) Address() string {
	return server.address
}

func (server *SimpleServer) IsAlive() bool {
	response, err := http.Get(server.Address())
	if err != nil {
		return false
	}
	if response.StatusCode == http.StatusOK {
		return true
	}
	return false
}

func (server *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	server.proxy.ServeHTTP(rw, r)
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		fmt.Printf("%s is Not Alive!\n", server.Address())
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server

}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.GetNextAvailableServer()
	fmt.Printf("Forwarding Request to Address: %s\n ", targetServer.Address())
	targetServer.Serve(rw, r)
}
