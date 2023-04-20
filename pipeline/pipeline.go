package pipeline

type Handler interface {
	// Handle continue if return True, else break
	Handle(passable any, next func(passable any))
}

type HandlerFunc func(p any, next func(passable any))

func (f HandlerFunc) Handle(p any, next func(passable any)) {
	f(p, next)
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

func (m *Pipeline) Then(handle func()) {
	fn := m.resolve(handle)
	fn(m.passable)
}

func (m *Pipeline) resolve(handle func()) func(passable any) {
	var fn = func(passable any) { handle() }
	for i := len(m.handler) - 1; i >= 0; i-- {
		fn = func(carry func(passable any), item Handler) func(passable any) {
			return func(passable any) {
				item.Handle(passable, carry)
			}
		}(fn, m.handler[i])
	}
	return fn
}
