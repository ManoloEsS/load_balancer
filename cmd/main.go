package main

import (
	"fmt"
	"log"
	"net/http"

	loadserver "github.com/ManoloEsS/load_balancer/internal/load_server"
)

type loadBalancer struct {
	current_server loadserver.Server
	http_server    *http.Server
}

func main() {

	servers := loadserver.CreateLoadServers(5)

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

func (lb *loadBalancer) balance(servers []loadserver.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		n := len(servers)

		lb.current_server.Channel <- r
		lb.current_server = servers[(lb.current_server.Address+1)%n]
	}

}

func (lb *loadBalancer) close(servers []loadserver.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, s := range servers {
			fmt.Printf("closing channel %d\n", s.Address)
			close(s.Channel)
		}

		log.Fatal()
	}
}
