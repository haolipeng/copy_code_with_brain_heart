package main

import (
	"log"
	"net/http"
)

type server struct {
}

func (p *server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//打印http的url请求
	log.Println(req.URL.Path)
	resp.Write([]byte("hello world"))
}

func main() {
	var s server
	addr := "localhost:9999"
	http.ListenAndServe(addr, &s)
}
