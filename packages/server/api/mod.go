package api

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/sovlookup/p2p/libp2p"
)

func Register(s *ghttp.Server, n *libp2p.Node) {
	RegisterPeerHandler(s.Group("/api"), n)
}

func RegisterPeerHandler(api *ghttp.RouterGroup, n *libp2p.Node) {
	api.GET("/", func(r *ghttp.Request) {
		r.Response.WriteJsonExit(n.Info())
	})

	api.GET("/peers", func(r *ghttp.Request) {
		r.Response.WriteJsonExit(n.Peers())
	})

	api.GET("/connect/*addr", func(r *ghttp.Request) {
		err := n.Connect("/" + r.Get("addr").String())
		if err != nil {
			if err.Error() == "failed to find peers: failed to find any peer in table" {
				r.Response.WriteStatus(404)
			} else {
				r.Response.WriteStatus(400)
			}
			r.Response.WriteOverExit(err.Error())
		}
	})
}
