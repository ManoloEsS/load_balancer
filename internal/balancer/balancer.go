package balancer

import (
	"github.com/ManoloEsS/load_balancer/internal/config"
	"github.com/ManoloEsS/load_balancer/internal/gateway"
	"github.com/ManoloEsS/load_balancer/internal/proxy"
	"github.com/ManoloEsS/load_balancer/internal/selector"
)

type LoadBalancer struct {
	gatewayServer *gateway.GatewayServer
}

func NewLoadBalancer(cfg *config.Config) (*LoadBalancer, error) {
	slctr := selector.NewSelector(cfg.Algorithm, cfg.Servers)

	lb := &LoadBalancer{
		gatewayServer: gateway.NewGatewayServer(proxy.NewProxy(slctr), cfg.ProxyAddress),
	}

	return lb, nil
}

func (lb *LoadBalancer) StartLoadBalancer() {
	lb.gatewayServer.Run()
}
