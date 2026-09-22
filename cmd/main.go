package main

import (
	"fmt"
	"log"
	"net/http"
)

type loadServer struct {
	address int
	channel chan *http.Request
}

type loadBalancer struct {
	current_server loadServer
	http_server    *http.Server
}

func main() {

	servers := createLoadServers(5)

	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lb := loadBalancer{
		http_server:    s,
		current_server: servers[0],
	}

	mux.HandleFunc("/", lb.balance(servers))
	mux.HandleFunc("/quit", lb.close(servers))

	fmt.Printf("Listening on port %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func (lb *loadBalancer) balance(servers []loadServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		n := len(servers)

		lb.current_server.channel <- r
		lb.current_server = servers[(lb.current_server.address+1)%n]
	}

}

func (lb *loadBalancer) close(servers []loadServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, s := range servers {
			fmt.Printf("closing channel %d\n", s.address)
			close(s.channel)
		}

		log.Fatal()
	}
}

func (ls *loadServer) serve() {
	go func() {
		for r := range ls.channel {
			fmt.Println(r.URL.Path)
			fmt.Printf("served by server: %d\n", ls.address)
		}
	}()
}

func createLoadServers(n int) []loadServer {
	servers := []loadServer{}
	for i := range n {
		srv := loadServer{
			address: i,
			channel: make(chan *http.Request),
		}

		servers = append(servers, srv)

		srv.serve()
	}
	return servers
}
