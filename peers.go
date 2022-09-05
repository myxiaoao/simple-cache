package simple_cache

import pb "simple-cache/cachepb"

// 在这之前已经实现了流程 ⑴ 和 ⑶，今天实现流程 ⑵，从远程节点获取缓存值。
// 我们进一步细化流程 ⑵：
// 使用一致性哈希选择节点        是                                   是
//    |-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
//                    |  否                                    ↓  否
//                    |----------------------------> 回退到本地节点处理。

// 在这里，抽象出 2 个接口，PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
// 接口 PeerGetter 的 Get() 方法用于从对应 group 查找缓存值。PeerGetter 就对应于上述流程中的 HTTP 客户端。

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
type PeerGetter interface {
	// Get Get(group string, key string) ([]byte, error)
	Get(in *pb.Request, out *pb.Response) error
}
