// +build windows

package raspberry

type P struct {
	pin     int
	edge    Edge
	handler func(*P)
}

func Open() error {
	return nil
}

func Close() error {
	return nil
}

func NewPin(p int) *P {
	return &(P{pin: p})
}

func (p *P) Watch(edge Edge, handler func(*P)) {
	p.handler = handler
	p.edge = edge
}

func (p *P) Unwatch() {
}

func (p *P) TestPin(edge Edge) {
	switch {
	case p.edge == EdgeNone, edge == EdgeNone:
		return

	case edge == EdgeBoth:
		// if edge is EdgeBoth, handler is called twice
		if p.edge == EdgeBoth {
			p.handler(p)
		}

		if p.edge == EdgeBoth || p.edge == EdgeFalling || p.edge == EdgeRising {
			p.handler(p)
		}
	case edge == EdgeFalling:
		if p.edge == EdgeBoth || p.edge == EdgeFalling {
			p.handler(p)
		}
	case edge == EdgeRising:
		if p.edge == EdgeBoth || p.edge == EdgeRising {
			p.handler(p)
		}
	}
}

func (p *P) Input() {
}

func (p *P) PullUp() {
}

func (p *P) PullDown() {
}

func (p *P) Pin() int {
	return p.pin
}
