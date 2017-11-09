package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("链接rpc服务器失败:", err)
	}
	var reply int
	err = client.Call("Watcher.GetIPCount", "self", &reply)
	if err != nil {
		fmt.Println("调用远程服务失败", err)
	}
	fmt.Println("远程服务返回结果：", reply)
}
