package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var sendChannel chan int = make(chan int, 1)

type DuplexUdp struct {
	Conn       *net.UDPConn
	ListenAddr *net.UDPAddr
	ServerAddr *net.UDPAddr
}

func (du *DuplexUdp) SetServerAddr(addr string) error {
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	du.ServerAddr = uaddr
	return err
}

func (du *DuplexUdp) SetListenAddr(addr string) error {
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	du.ListenAddr = uaddr
	return err
}

func (du *DuplexUdp) DialUDP() error {
	conn, err := net.DialUDP("udp", du.ListenAddr, du.ServerAddr)
	du.Conn = conn
	return err
}

func (du DuplexUdp) Recv() {
	data := make([]byte, 2048)
	for {
		n, raddr, err := du.Conn.ReadFromUDP(data)
		CheckErrorOnExit(err)
		fmt.Println("message:", string(data[:n]))
		fmt.Println("remote addr:", raddr)
		splitData := strings.Split(string(data[:n]), ":")
		ret_data, err := strconv.Atoi(splitData[1])
		CheckErrorOnExit(err)
		if splitData[0] == "client" {
			du.Send([]byte("server:" + strconv.Itoa(ret_data+1)))
		} else {
			sendChannel <- ret_data
		}
	}
}

func (du DuplexUdp) Send(data []byte) {
	_, err := du.Conn.Write(data)
	CheckErrorOnExit(err)
}

func main() {
	listenAddr := "127.0.0.1:8989"
	serverAddr := "127.0.0.1:9898"

	duplex := DuplexUdp{}
	err := duplex.SetServerAddr(serverAddr)
	CheckErrorOnExit(err)
	err = duplex.SetListenAddr(listenAddr)
	CheckErrorOnExit(err)

	err = duplex.DialUDP()
	CheckErrorOnExit(err)

	go duplex.Recv()

	time.Sleep(time.Second * 3)

	for i := 0; ; i++ {
		duplex.Send([]byte("client:" + strconv.Itoa(i)))
		select {
		case ret_data := <-sendChannel:
			if ret_data == i+1 {
				fmt.Println("send success!")
			} else {
				fmt.Println("send failed!")
			}
		case <-time.After(time.Second * 2):
			fmt.Println("send failed!")
		}
		time.Sleep(time.Second * 2)
	}

	time.Sleep(time.Second * 3)
}

func CheckErrorOnExit(err error) {
	if err != nil {
		panic(err)
	}
}
