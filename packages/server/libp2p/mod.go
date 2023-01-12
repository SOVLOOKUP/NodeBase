package libp2p

import (
	"context"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	webtransport "github.com/libp2p/go-libp2p/p2p/transport/webtransport"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	h   host.Host
	cfg *NodeOpt
}

type Info struct {
	ID       string   `json:"id"`
	CertHash []string `json:"certHash"`
}

// 连接 Peer
func (n *Node) Connect(addr string) error {
	peerAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	addrinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
	if err != nil {
		return err
	}
	lc, lcCancel := context.WithTimeout(n.cfg.Ctx, n.cfg.ConnectTimeout)
	defer lcCancel()
	return n.h.Connect(lc, *addrinfo)
}

// 获取 Peers
func (n *Node) Peers() peer.IDSlice {
	return n.h.Network().Peers()
}

// 根据 Peer ID 获取 Name
// todo 分布式储存
func (n *Node) PeerName(id string) (string, error) {
	peer, err := peer.Decode(id)
	if err != nil {
		return "", err
	}
	name, err := n.h.Peerstore().Get(peer, "name")
	if err != nil {
		return "", err
	}
	return name.(string), nil
}

// 获取当前节点监听的 UDP 端口号
func (n *Node) Port() string {
	return n.cfg.Port
}

// 获取 Name ID CertHash 等当前节点信息
func (n *Node) Info() Info {
	addrs := n.h.Addrs()

	i := Info{
		ID: n.h.ID().String(),
	}

	res := slice.Filter(addrs, func(_ int, addr multiaddr.Multiaddr) bool {
		return strings.Contains(addr.String(), "webtransport")
	})

	if len(res) > 0 {
		i.CertHash = slice.Map(
			multiaddr.FilterAddrs(multiaddr.Split(res[0]),
				func(addr multiaddr.Multiaddr) bool {
					return strings.Contains(addr.String(), "certhash")
				}),
			func(_ int, addr multiaddr.Multiaddr) string {
				return strings.TrimPrefix(addr.String(), "/certhash/")
			})

	} else {
		i.CertHash = []string{}
	}

	return i
}

// 创建当前节点
func New(opt *NodeOpt) (Node, error) {
	cfg := []config.Option{
		libp2p.Transport(quic.NewTransport),
		libp2p.Transport(webtransport.New),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/"+opt.Port+"/quic-v1",
			"/ip4/0.0.0.0/udp/"+opt.Port+"/quic-v1/webtransport",
			"/ip6/::/udp/"+opt.Port+"/quic-v1",
			"/ip6/::/udp/"+opt.Port+"/quic-v1/webtransport",
		),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err := dht.New(opt.Ctx, h)

			return idht, err
		}),
	}

	// 启动节点
	h, err := libp2p.New(append(cfg, opt.Libp2pOpts...)...)
	if err != nil {
		return Node{}, err
	}

	// 回写使用的端口号
	addrs := h.Addrs()
	if len(addrs) > 0 {
		opt.Port = strings.Split(multiaddr.Split(addrs[0])[1].String(), "/")[2]
	}

	return Node{
		h,
		opt,
	}, nil
}
