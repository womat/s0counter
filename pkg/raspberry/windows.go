// +build windows

package raspberry

type Line struct {
	pin     int
	edge    Edge
	handler func(*Line)
}

type Chip struct {
}

func Open() (*Chip, error) {
	return &Chip{}, nil
}

func (c *Chip) Close() error {
	return nil
}

func (c *Chip) NewPin(p int) *Line {
	return &(Line{pin: p})
}

func (l *Line) Watch(edge Edge, handler func(*Line)) error {
	l.handler = handler
	l.edge = edge
	return nil
}

func (l *Line) Unwatch() {
}

func (l *Line) TestPin(edge Edge) {
	switch {
	case l.edge == EdgeNone, edge == EdgeNone:
		return

	case edge == EdgeBoth:
		// if edge is EdgeBoth, handler is called twice
		if l.edge == EdgeBoth {
			l.handler(l)
		}

		if l.edge == EdgeBoth || l.edge == EdgeFalling || l.edge == EdgeRising {
			l.handler(l)
		}
	case edge == EdgeFalling:
		if l.edge == EdgeBoth || l.edge == EdgeFalling {
			l.handler(l)
		}
	case edge == EdgeRising:
		if l.edge == EdgeBoth || l.edge == EdgeRising {
			l.handler(l)
		}
	}
}

func (l *Line) Input() {
}

func (l *Line) PullUp() {
}

func (l *Line) PullDown() {
}

func (l *Line) Pin() int {
	return l.pin
}

func (l *Line) Read() bool {
	return true
}
