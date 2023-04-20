package pipeline

type Next func(passable any) any

type Callback func() any

type Handler interface {
	// Handle continue if return True, else break
	Handle(passable any, next Next) any
}

type HandlerFunc func(p any, next Next) any

func (f HandlerFunc) Handle(p any, next Next) any {
	return f(p, next)
}

type Pipeline struct {
	handler  []Handler
	passable any
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (m *Pipeline) Send(passable any) *Pipeline {
	m.passable = passable
	return m
}

func (m *Pipeline) Through(hds []Handler) *Pipeline {
	m.handler = hds
	return m
}

func (m *Pipeline) Then(handle Callback) any {
	fn := m.resolve(handle)
	return fn(m.passable)
}

func (m *Pipeline) resolve(handle Callback) Next {
	var fn = func(passable any) any { return handle() }
	for i := len(m.handler) - 1; i >= 0; i-- {
		fn = func(carry Next, item Handler) Next {
			return func(passable any) any {
				return item.Handle(passable, carry)
			}
		}(fn, m.handler[i])
	}
	return fn
}
