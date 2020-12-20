// +build windows

package raspberry

import (
	"fmt"
	"time"
)

type Line struct {
	pin  int
	edge Edge
	// the bounceTime defines the key bounce time (ms)
	// the value 0 ignores key bouncing
	bounceTime time.Duration
	// while bounceTimer is running, new signal are ignored (suppress key bouncing)
	bounceTimer *time.Timer
	lastLevel   bool
	handler     func(*Line)
}

type Chip struct {
}

func Open() (*Chip, error) {
	lines = map[int]*Line{}
	return &Chip{}, nil
}

func (c *Chip) Close() error {
	return nil
}

func (c *Chip) NewPin(p int) (*Line, error) {
	if _, ok := lines[p]; ok {
		return nil, fmt.Errorf("pin %v already used", p)
	}

	lines[p] = &Line{pin: p, bounceTimer: time.NewTimer(0)}
	return lines[p], nil
}

func (l *Line) Watch(edge Edge, handler func(*Line)) error {
	l.handler = handler
	l.edge = edge
	return nil
}

func (l *Line) SetBounceTime(t time.Duration) *Line {
	l.bounceTime = t
	return l
}

func (l *Line) BounceTime() time.Duration {
	return l.bounceTime
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
	// check if map with pin struct exists
	l, ok := lines[pin.Pin()]
	if !ok {
		return
	}

	// if debounce is inactive, call handler function and returns
	if l.bounceTime == 0 {
		l.lastLevel = l.Read()
		l.handler(pin)
		return
	}

	select {
	case <-l.bounceTimer.C:
		// if bounce Timer is expired, accept new signals
		l.bounceTimer.Reset(l.bounceTime)
	default:
		// if bounce Timer is still running, ignore single
		return
	}

	go func(l *Line) {
		// wait until bounce Timer is expired and check if the pin has still the correct level
		// the correct level depends on the edge configuration
		<-l.bounceTimer.C
		l.bounceTimer.Reset(0)

		switch l.edge {
		case EdgeBoth:
			if l.Read() != l.lastLevel {
				l.lastLevel = l.Read()
				l.handler(pin)
			}
		case EdgeFalling:
			if !l.Read() {
				l.lastLevel = l.Read()
				l.handler(pin)
			}
		case EdgeRising:
			if l.Read() {
				l.lastLevel = l.Read()
				l.handler(pin)
			}
		}
	}(l)
}
