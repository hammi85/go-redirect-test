package server

// Context holds context
type Context struct {
	state         *State
	parentContext *Context
	complete      bool
}

// NewContext creates a new app context
func NewContext(state *State, parentCtx *Context) *Context {
	c := &Context{state: state, parentContext: parentCtx}

	if parentCtx != nil {
		c.complete = parentCtx.complete
	}

	return c
}
