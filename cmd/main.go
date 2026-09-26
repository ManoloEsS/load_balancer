package main

import (
	"fmt"
	"log"

	"github.com/ManoloEsS/load_balancer/internal/balancer"
	"github.com/ManoloEsS/load_balancer/internal/config"
)

func main() {

	// load configuration
	config, err := config.LoadConfig(config.Config_file)
	if err != nil {
		fmt.Printf("config: %w", err)
	}

	// create load balancer
	lb, err := balancer.NewLoadBalancer(config)
	if err != nil {
		log.Fatalf("could not start load balancer: %v", err)
	}

	lb.StartLoadBalancer()
}
