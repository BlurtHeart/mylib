package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Watcher int

func (w *Watcher) GetIPCount(arg string, result *int) error {
	fp, err := os.Open("filelist.conf")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var wg sync.WaitGroup
	var lock sync.Mutex
	var total int
	br := bufio.NewReader(fp)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		wg.Add(1)
		go func() {
			cmd := "grep " + arg + " " + string(a)
			count := handleCmd(cmd)
			lock.Lock()
			total += count
			lock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	*result = total
	return nil
}

func main() {
	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("监听失败，端口可能已经被占用")
	}
	fmt.Println("正在监听12345端口")
	http.Serve(l, nil)
}

func handleCmd(a string) int {
	fmt.Println("cmd:", a)
	cmd := exec.Command("sh", "-c", a)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	result := strings.Split(string(out), "\n")
	fmt.Println("result:", result, "len:", len(result))
	if len(result) != 0 {
		return len(result) - 1
	}
	return 0
}
