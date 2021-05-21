package match

import (
	"sync"
)

// Context is events passed down to cards, allowing them to perform actions
// without having a direct reference to the match, players etc
type Context struct {
	Match   *Match
	Event   interface{}
	cancel  bool
	preFxs  []func()
	mainFx  func()
	postFxs []func()

	mutex *sync.Mutex
}

// HandlerFunc is a function with a match context as argument
type HandlerFunc func(card *Card, ctx *Context)

// NewContext returns a new match context
func NewContext(m *Match, e interface{}) *Context {

	ctx := &Context{
		Match: m,
		Event: e,
		mutex: &sync.Mutex{},
	}

	return ctx
}

// ScheduleBefore allows you to run the logic before the main logic,
func (c *Context) ScheduleBefore(handlers ...func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.preFxs = append(c.preFxs, handlers...)
}

// Override holds or overrides the main logic,
func (c *Context) Override(handler func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.mainFx = handler
}

// ScheduleAfter allows you to run the logic at the end of the context flow,
// after the default behaviour
func (c *Context) ScheduleAfter(handlers ...func()) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.postFxs = append(c.postFxs, handlers...)
}

// InterruptFlow stops the context flow, cancelling the default behaviour
func (c *Context) InterruptFlow() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cancel = true
}

// Cancelled returns whether or not the context was cancelled
func (c *Context) Cancelled() bool {
	return c.cancel
}
