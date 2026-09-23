package loadserver

import (
	"fmt"
	"net/http"
)

type Server struct {
	Address int
	Channel chan *http.Request
}

func (ls *Server) serve() {
	go func() {
		for r := range ls.Channel {
			fmt.Println(r.URL.Path)
			fmt.Printf("served by server: %d\n", ls.Address)
		}
	}()
}

func CreateLoadServers(n int) []Server {
	servers := []Server{}
	for i := range n {
		srv := Server{
			Address: i,
			Channel: make(chan *http.Request),
		}

		servers = append(servers, srv)

		srv.serve()
	}
	return servers
}
