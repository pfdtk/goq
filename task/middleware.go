package task

type Middleware interface {
	// Handle continue if return True, else break
	Handle(p *Passable, next func(passable *Passable))
}

type FuncHandler func(p *Passable, next func(passable *Passable))

func (f FuncHandler) Handle(p *Passable, next func(passable *Passable)) {
	f(p, next)
}

type Passable struct {
	t Task
	j *Job
}

func NewPassable(t Task, j *Job) *Passable {
	return &Passable{t: t, j: j}
}

type MiddlewarePipeline struct {
	middleware []Middleware
	passable   *Passable
	abort      bool
}

func NewMiddlewarePipeline() *MiddlewarePipeline {
	return &MiddlewarePipeline{}
}

func (m *MiddlewarePipeline) Send(passable *Passable) *MiddlewarePipeline {
	m.passable = passable
	return m
}

func (m *MiddlewarePipeline) Through(md []Middleware) *MiddlewarePipeline {
	m.middleware = md
	return m
}

func (m *MiddlewarePipeline) Then(handle func()) {
	fn := m.resolve(handle)
	fn(m.passable)
}

func (m *MiddlewarePipeline) resolve(handle func()) func(passable *Passable) {
	var fn = func(passable *Passable) { handle() }
	for i := len(m.middleware) - 1; i >= 0; i-- {
		fn = func(carry func(passable *Passable), item Middleware) func(passable *Passable) {
			return func(passable *Passable) {
				item.Handle(passable, carry)
			}
		}(fn, m.middleware[i])
	}
	return fn
}
