package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blurtheart/mylib/log"
	"context"
)

func main() {
	http.HandleFunc("/", log.Decorator(handler))
	panic(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, int(42), int64(100))

	log.Println(ctx, "handler started")
	defer log.Println(ctx, "handler ended")

	select {
	case <-time.After(time.Second * 5):
		fmt.Fprint(w, "hello world")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
