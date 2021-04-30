package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

//HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	//this peer's base URL,"https://example.net:8000"
	self     string //记录自己的地址，包括主机名/ip和端口
	basePath string //节点间通信的前缀，默认是/_geecache/，主机上可能还承载其他业务，加一段Path是一个好习惯
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
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
