// +build windows

package raspberry

import (
	"fmt"
	"time"
)

type Line struct {
	pin           int
	edge          Edge
	debounceTime  time.Duration
	debounceTimer *time.Timer
	lastLevel     bool
	handler       func(*Line)
}

type Chip struct {
}

func Open() (*Chip, error) {
	lines = []*Line{}
	return &Chip{}, nil
}

func (c *Chip) Close() error {
	return nil
}

func (c *Chip) NewPin(p int) (*Line, error) {
	for _, l := range lines {
		if l.pin == p {
			return nil, fmt.Errorf("pin %v already used", p)
		}
	}

	l := Line{pin: p, debounceTimer: time.NewTimer(0)}
	lines = append(lines, &l)
	return &l, nil
}

func (l *Line) Watch(edge Edge, handler func(*Line)) error {
	l.handler = handler
	l.edge = edge
	return nil
}

func (l *Line) SetDebounceTimer(t time.Duration) *Line {
	l.debounceTime = t
	return l
}

func (l *Line) DebounceTimer() time.Duration {
	return l.debounceTime
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
			handler(l)
		}

		if l.edge == EdgeBoth || l.edge == EdgeFalling || l.edge == EdgeRising {
			handler(l)
		}
	case edge == EdgeFalling:
		if l.edge == EdgeBoth || l.edge == EdgeFalling {
			handler(l)
		}
	case edge == EdgeRising:
		if l.edge == EdgeBoth || l.edge == EdgeRising {
			handler(l)
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
	return false
}

func handler(pin *Line) {
	for _, l := range lines {
		if l == pin {
			if l.debounceTime == 0 {
				l.lastLevel = l.Read()
				l.handler(pin)
				return
			}

			select {
			case <-l.debounceTimer.C:
				l.debounceTimer.Reset(l.debounceTime)
			default:
				return
			}

			go func(p *Line, l Line) {
				<-l.debounceTimer.C
				l.debounceTimer.Reset(0)

				switch l.edge {
				case EdgeBoth:
					if p.Read() != l.lastLevel {
						l.lastLevel = l.Read()
						l.handler(pin)
					}
				case EdgeFalling:
					if !p.Read() {
						l.lastLevel = l.Read()
						l.handler(pin)
					}
				case EdgeRising:
					if p.Read() {
						l.lastLevel = l.Read()
						l.handler(pin)
					}
				}
			}(pin, *l)
			return
		}
	}
}
