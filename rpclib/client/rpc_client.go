package main

import (
	"log"
	"net/rpc"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("give a string to count")
	}
	arg := os.Args[1]
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalln("链接rpc服务器失败:", err)
	}
	var reply int
	err = client.Call("Watcher.GetCount", arg, &reply)
	if err != nil {
		log.Fatalln("调用远程服务失败", err)
	}
	log.Println("远程服务返回结果：", reply)
}
