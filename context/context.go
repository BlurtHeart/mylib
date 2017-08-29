// implement standard library's context

package context

import (
	"errors"
	"fmt"
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
	propagateCancel(parent, c)
	return c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) *CancelContext {
	return &CancelContext{
		Context: parent,
		done:    make(chan struct{}),
	}
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

func propagateCancel(parent Context, child canceler) {
	if parent.Done() == nil {
		return // parent is never canceled
	}
	if p, ok := parentCancelCtx(parent); ok {
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]bool)
			}
			p.children[child] = true
		}
	} else {
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

func parentCancelCtx(parent Context) (*CancelContext, bool) {
	for {
		switch c := parent.(type) {
		case *CancelContext:
			return c, true
		default:
			return nil, false
		}
	}
}

type CancelContext struct {
	Context
	done     chan struct{}
	err      error
	children map[canceler]bool
}

func (cctx *CancelContext) String() string {
	return fmt.Sprintf("%v.WithCancel", cctx.Context)
}

func (cctx *CancelContext) Done() <-chan struct{} {
	return cctx.done
}

func (cctx *CancelContext) Err() error {
	return cctx.err
}

func (cctx *CancelContext) cancel(removeFromParent bool, err error) {
	cctx.err = err
	close(cctx.done)
	return
}
