package match

// Context is events passed down to cards, allowing them to perform actions
// without having a direct reference to the match, players etc
type Context struct {
	match   *Match
	event   interface{}
	cancel  bool
	preFxs  []func()
	mainFx  func()
	postFxs []func()
}

// HandlerFunc is a function with a match context as argument
type HandlerFunc func(card *Card, ctx *Context)

// NewContext returns a new match context
func NewContext(m *Match, e interface{}) *Context {

	ctx := &Context{
		match: m,
		event: e,
	}

	return ctx
}

// Match ...
func (ctx *Context) Match() *Match {
	return ctx.match
}

// Event ...
func (ctx *Context) Event() interface{} {
	return ctx.event
}

// ScheduleBefore allows you to run the logic before the main logic,
func (ctx *Context) ScheduleBefore(handlers ...func()) {
	ctx.preFxs = append(ctx.preFxs, handlers...)
}

// Override holds or overrides the main logic,
func (ctx *Context) Override(handler func()) {
	ctx.mainFx = handler
}

// ScheduleAfter allows you to run the logic at the end of the context flow,
// after the default behaviour
func (ctx *Context) ScheduleAfter(handlers ...func()) {
	ctx.postFxs = append(ctx.postFxs, handlers...)
}

// InterruptFlow stops the context flow, cancelling the default behaviour
func (ctx *Context) InterruptFlow() {
	ctx.cancel = true
}

// Cancelled returns whether or not the context was cancelled
func (ctx *Context) Cancelled() bool {
	return ctx.cancel
}
