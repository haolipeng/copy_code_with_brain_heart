本节课按照如下流程来讲
首先，讲解下一个普通的http服务器如何编写
然后，讲解下分布式缓存中的通信流程，为什么会有前缀信息
我们约定访问路径的格式为:
http://example.com/api/xxxx
http://example.com/<basepath>/<groupname>/<key>
basepath:作为节点间通信地址的前缀，比如/_geecache/的url地址是用于缓存请求业务的，
而/api/是用于对外接口api业务的。

然后，定义分布式缓存中HTTP通信的核心结构体HTTPPool，如何定义

最后，总结一下实现接口函数ServeHTTP的步骤
1、先校验请求合法性，即url中是否有basepath字符串
2、解析出url路径中的groupname字段
3、解析出url路径中的key字段
4、找到groupname对应的group对象，根据缓存key在group对象中获取到相应的缓存value数值
5、将获取的缓存值写回到http响应中