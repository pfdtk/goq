package task

type Middleware interface {
	// Handle continue if return True, else break
	Handle(p *Passable, next func(passable *Passable))
}

type MiddlewareFunc func(p *Passable, next func(passable *Passable))

func (f MiddlewareFunc) Handle(p *Passable, next func(passable *Passable)) {
	f(p, next)
}

type Passable struct {
	t Task
	j *Job
}

func NewPassable(t Task, j *Job) *Passable {
	return &Passable{t: t, j: j}
}

type Pipeline struct {
	middleware []Middleware
	passable   *Passable
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (m *Pipeline) Send(passable *Passable) *Pipeline {
	m.passable = passable
	return m
}

func (m *Pipeline) Through(md []Middleware) *Pipeline {
	m.middleware = md
	return m
}

func (m *Pipeline) Then(handle func()) {
	fn := m.resolve(handle)
	fn(m.passable)
}

func (m *Pipeline) resolve(handle func()) func(passable *Passable) {
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
