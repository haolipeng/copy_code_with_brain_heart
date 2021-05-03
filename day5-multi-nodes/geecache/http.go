package geecache

import (
	"copy-group-cache/day5-multi-nodes/geecache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_geecache/"
const defaultReplicas = 50

//HTTPPool实现了两种能力：
//1.提供http服务的能力
//2.根据指定的key，创建http客户端从远程节点获取缓存值的能力

type httpGetter struct {
	baseURL string //要访问的远程节点的地址
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	//构造http请求
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned:%v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

//HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	//this peer's base URL,"https://example.net:8000"
	self        string                 //记录自己的地址，包括主机名/ip和端口
	basePath    string                 //节点间通信的前缀，默认是/_geecache/，主机上可能还承载其他业务，加一段Path是一个好习惯
	mu          sync.Mutex             //并发访问控制
	peers       *consistenthash.Map    //用来根据指定的key来返回缓存节点
	httpGetters map[string]*httpGetter //维护远程节点和对应的HttpGetter的映射关系，每个远程节点对应一个httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Set 此函数不可重复调用
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//1.创建一致性哈希存储缓存节点
	//2.添加缓存节点服务器
	//3.每个缓存节点对应一个httpGetter
	p.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

// PickPeer 根据具体的key返回对应的缓存
func (p *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick Peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//ServeHTTP的http请求访问格式为 /<basepath>/<groupname>/<key>
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//1.判断访问路径的前缀是否是basePath
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	//打印出http的url路径和方法
	p.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	// 2.从请求url中解析出groupname和key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 { //单元测试出这里有问题
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//group name
	groupName := parts[0]
	key := parts[1]

	//通过groupName得到group实例
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	//通过key得到缓存数据
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//将缓存数据作为结果写入到http响应中
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
