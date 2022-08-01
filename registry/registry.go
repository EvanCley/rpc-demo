package registry

type Registry interface {
	Register(option RegisterOption, provider ...Provider)   // 注册
	Unregister(option RegisterOption, provider ...Provider) // 注销
	GetServiceList() []Provider                             // 获取服务列表
	Watch() Watcher                                         // 监听服务列表的变化
	Unwatch(watcher Watcher)                                // 取消监听
}

type RegisterOption struct {
	AppKey string // AppKey用于唯一标识某个应用，比如com.meituan.demo.rpc.server
}

type Watcher interface {
	Next() (*Event, error) // 获取下一次服务列表的更新
	Close()
}

type EventAction byte

const (
	Create EventAction = iota
	Update
	Delete
)

// Event 事件表示一次更新
type Event struct {
	AppKey    string
	Providers []Provider
}

// Provider 某个具体的服务提供者
type Provider struct {
	ProviderKey string // Network+"@"+Addr
	Network     string
	Addr        string
	Meta        map[string]interface{}
}

// Peer2PeerDiscovery 下面就是服务端与客户端直连的接口实现
type Peer2PeerDiscovery struct {
	providers []Provider
}

func (p *Peer2PeerDiscovery) Register(option RegisterOption, providers ...Provider) {
	p.providers = providers
}

func (p *Peer2PeerDiscovery) Unregister(option RegisterOption, provider ...Provider) {
	p.providers = []Provider{}
}

func (p *Peer2PeerDiscovery) GetServiceList() []Provider {
	return p.providers
}

func (p *Peer2PeerDiscovery) Watch() Watcher {
	return nil
}

func (p *Peer2PeerDiscovery) Unwatch(watcher Watcher) {
	return
}

func (p *Peer2PeerDiscovery) WithProvider(provider Provider) *Peer2PeerDiscovery {
	p.providers = append(p.providers, provider)
	return p
}

func (p *Peer2PeerDiscovery) WithProviders(providers []Provider) *Peer2PeerDiscovery {
	for _, provider := range providers {
		p.providers = append(p.providers, provider)
	}
	return p
}

func NewPeer2PeerRegistry() *Peer2PeerDiscovery {
	r := &Peer2PeerDiscovery{}
	return r
}
