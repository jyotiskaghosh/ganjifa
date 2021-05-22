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
func (ctx *Context) ScheduleBefore(handlers ...func()) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.preFxs = append(ctx.preFxs, handlers...)
}

// Override holds or overrides the main logic,
func (ctx *Context) Override(handler func()) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.mainFx = handler
}

// ScheduleAfter allows you to run the logic at the end of the context flow,
// after the default behaviour
func (ctx *Context) ScheduleAfter(handlers ...func()) {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.postFxs = append(ctx.postFxs, handlers...)
}

// InterruptFlow stops the context flow, cancelling the default behaviour
func (ctx *Context) InterruptFlow() {

	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.cancel = true
}

// Cancelled returns whether or not the context was cancelled
func (ctx *Context) Cancelled() bool {
	return ctx.cancel
}
