package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Watcher int

func (w *Watcher) GetCount(arg string, result *int) error {
	filelist := []string{"111", "222"}

	var wg sync.WaitGroup
	var lock sync.Mutex
	var total int
	for _, filename := range filelist {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			cmd := "grep " + arg + " " + string(file) + "|wc -l"
			countStr, err := handleCmd(cmd)
			if err != nil {
				log.Println(err)
				return
			}
			count, err := strconv.Atoi(strings.TrimSpace(countStr))
			if err != nil {
				log.Println(err)
				return
			}
			lock.Lock()
			total += count
			lock.Unlock()
		}(filename)
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
		log.Fatalf("监听失败，端口可能已经被占用")
	}
	log.Println("正在监听12345端口")
	http.Serve(l, nil)
}

func handleCmd(a string) (string, error) {
	cmd := exec.Command("sh", "-c", a)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
