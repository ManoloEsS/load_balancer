package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/ManoloEsS/load_balancer/internal/selector"
)

// type Proxy struct {
// 	AlgoFunc func(r *httputil.ProxyRequest)
//
// }

func NewProxy(selector selector.BackendSelector) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			next := selector.Select()
			nextUrl, err := url.Parse("http://" + next.Address)
			if err != nil {
				return
			}
			r.SetXForwarded()
			r.SetURL(nextUrl)
		},
	}

	return proxy
}

// func (rp *Proxy) SetRewriteFunc(rwfunc func(*httputil.ProxyRequest)) {
// 	rp.AlgoFunc = rwfunc
// }
