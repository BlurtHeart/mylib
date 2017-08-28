// implement standard library's context

package context

import (
	"errors"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyCtx int

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyCtx) Done() <-chan struct{} {
	return nil
}

func (emptyCtx) Err() error {
	return nil
}

func (emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (ctx *emptyCtx) String() string {
	switch ctx {
	case background:
		return "context.Background"
	case todo:
		return "context.TODO"
	}
	return "Unkown empty context"
}

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context {
	return background
}

func TODO() Context {
	return todo
}

type cancelFunc func()

var Canceled = errors.New("context canceled")

func WithCancel(parent Context) (Context, cancelFunc) {
	c := newCancelCtx(parent)
	return c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) *CancelContext {
	return &CancelContext{
		Context: parent,
		done:    make(chan struct{}),
	}
}

type CancelContext struct {
	Context
	done chan struct{}
	err  error
}

func (cctx *CancelContext) String() string {
	return "context.Background.WithCancel"
}

func (cctx *CancelContext) Done() <-chan struct{} {
	return cctx.done
}

func (cctx *CancelContext) cancel(removeFromParent bool, err error) {
	return
}
