// implement standard library's context

package context

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
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

type CancelFunc func()

var Canceled = errors.New("context canceled")

var DeadlineExceeded = errors.New("context deadline exceeded")

func WithCancel(parent Context) (Context, CancelFunc) {
	c := newCancelCtx(parent)
	propagateCancel(parent, c)
	return c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) *cancelCtx {
	return &cancelCtx{
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
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]bool)
			}
			p.children[child] = true
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			// block here, wait for channel close of parent or child
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		case *cancelCtx:
			return c, true
		case *timerCtx:
			return c.cancelCtx, true
		case *valueCtx:
			parent = c.Context
		default:
			return nil, false
		}
	}
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.children != nil {
		delete(p.children, child)
	}
}

type cancelCtx struct {
	Context
	done     chan struct{}
	err      error
	children map[canceler]bool
	mu       sync.Mutex
}

func (cctx *cancelCtx) String() string {
	return fmt.Sprintf("%v.WithCancel", cctx.Context)
}

func (cctx *cancelCtx) Done() <-chan struct{} {
	return cctx.done
}

func (cctx *cancelCtx) Err() error {
	cctx.mu.Lock()
	defer cctx.mu.Unlock()
	return cctx.err
}

func (cctx *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	cctx.mu.Lock()
	if cctx.err != nil {
		cctx.mu.Unlock()
		return // already canceled
	}
	cctx.err = err
	close(cctx.done)

	for k := range cctx.children {
		k.cancel(false, err)
	}
	cctx.children = nil
	cctx.mu.Unlock()

	if removeFromParent {
		removeChild(cctx.Context, cctx)
	}
	return
}

type timerCtx struct {
	*cancelCtx
	timer    *time.Timer
	deadline time.Time
}

func (c *timerCtx) String() string {
	return fmt.Sprintf("%v.WithDeadline(%s [%s])", c.cancelCtx.Context, c.deadline, c.deadline.Sub(time.Now()))
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  deadline,
	}
	propagateCancel(parent, c)
	d := deadline.Sub(time.Now())
	if d <= 0 {
		// deadline already passed
		c.cancel(true, DeadlineExceeded)
		return c, func() { c.cancel(true, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(d, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}

	return c, func() { c.cancel(true, Canceled) }
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithValue(parent Context, key, value interface{}) Context {
	c := &valueCtx{
		Context: parent,
		key:     key,
		value:   value,
	}
	return c
}

type valueCtx struct {
	Context
	key, value interface{}
}

func (c *valueCtx) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.value)
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.value
	}
	return c.Context.Value(key)
}
