package pipeline

type Next func(p any) any

type Handler func(p any, next func(p any) any) any

type Pipeline struct {
	handler  []Handler
	passable any
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (m *Pipeline) Send(p any) *Pipeline {
	m.passable = p
	return m
}

func (m *Pipeline) Through(hds ...Handler) *Pipeline {
	m.handler = hds
	return m
}

func (m *Pipeline) Then(handle Next) any {
	fn := m.resolve(handle)
	return fn(m.passable)
}

func (m *Pipeline) resolve(handle Next) Next {
	var fn = func(p any) any { return handle(p) }
	for i := len(m.handler) - 1; i >= 0; i-- {
		fn = func(carry Next, h Handler) Next {
			return func(p any) any {
				return h(p, carry)
			}
		}(fn, m.handler[i])
	}
	return fn
}
