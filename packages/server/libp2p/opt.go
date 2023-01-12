package libp2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

type NodeOpt struct {
	Ctx            context.Context
	Port           string
	ConnectTimeout time.Duration
	Libp2pOpts     []config.Option
}

// 启用自动中继，必须至少满足以下条件之一：
//
// 1. 提供静态中继节点 Addr
//
// 2. 提供动态 WithPeerSource 通道函数
func (o *NodeOpt) WithAutoRelay(opts ...autorelay.Option) {
	o.Libp2pOpts = append(o.Libp2pOpts, libp2p.EnableAutoRelay(opts...))
}

// 初始化默认配置
func DefaultOpt() *NodeOpt {
	return &NodeOpt{
		Ctx:            context.Background(),
		ConnectTimeout: time.Second * 15,
		Port:           "0",
	}
}
