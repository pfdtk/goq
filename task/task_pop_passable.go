package task

type PopPassable struct {
	callback func()
}

func NewPopPassable() *PopPassable {
	return &PopPassable{}
}

func (p *PopPassable) SetCallback(fn func()) {
	p.callback = fn
}

func (p *PopPassable) GetCallback() func() {
	return p.callback
}

func (p *PopPassable) ExecCallback() {
	if p.callback != nil {
		p.callback()
	}
}
