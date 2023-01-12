package main

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/sovlookup/p2p/api"
	"github.com/sovlookup/p2p/libp2p"
	_ "github.com/sovlookup/p2p/packed"
)

func main() {
	// 启动本地 P2P 节点服务
	opts := libp2p.DefaultOpt()
	n, err := libp2p.New(opts)
	if err != nil {
		panic(err)
	}

	// HTTP 服务
	println(">>>>>>>>>>> Listening on http://localhost:" + opts.Port)
	s := g.Server()
	// 设置前端资源
	s.SetServerRoot("dist")
	// 自动获取并设置 port
	port, err := convertor.ToInt(n.Port())
	if err != nil {
		panic(err)
	}
	s.SetPort(int(port))
	// 注册 API
	api.Register(s, &n)
	s.Run()
}
