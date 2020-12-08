// +build windows

package rpi

type Rpi struct {
}

type PinEmu struct {
	pin     int
	edge    Edge
	handler func(*PinEmu)
}

func Open() (*Rpi, error) {
	return &Rpi{}, nil
}

func (rpi *Rpi) Close() {
}

func (rpi *Rpi) NewPin(p int) *PinEmu {
	return &(PinEmu{pin: p})
}

func (p *PinEmu) Watch(edge Edge, handler interface{}) {
	if tmp, ok := handler.(func(*PinEmu)); ok {
		p.handler = tmp
		p.edge = edge
	}
}

func (p *PinEmu) Unwatch() {
}

func (p *PinEmu) TestPin(edge Edge) {
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

func (p *PinEmu) Input() {
}

func (p *PinEmu) PullUp() {
}

func (p *PinEmu) PullDown() {
}

func (p *PinEmu) Pin() int {
	return p.pin
}
