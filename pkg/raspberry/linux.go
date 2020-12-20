// +build !windows

package raspberry

import (
	"fmt"
	"time"

	"github.com/warthog618/gpio"
)

type Line struct {
	gpioPin       *gpio.Pin
	handler       func(*gpio.Pin)
	debounceTime  time.Duration
	debounceTimer *time.Timer
	edge          Edge
	lastLevel     gpio.Level
}
type Chip struct {
}

func Open() (*Chip, error) {
	if err := gpio.Open(); err != nil {
		return nil, err
	}

	lines = []*Line{}
	return &Chip{}, nil
}

func (c *Chip) Close() (err error) {
	return gpio.Close()
}

func (c *Chip) NewPin(p int) (*Line, error) {
	for _, l := range lines {
		if l.gpioPin.Pin() == p {
			return nil, fmt.Errorf("pin %v already used", p)
		}
	}

	l := Line{gpioPin: gpio.NewPin(p), debounceTimer: time.NewTimer(0)}
	lines = append(lines, &l)
	return &l, nil
}

func (l *Line) SetDebounceTimer(t time.Duration) *Line {
	l.debounceTime = t
	return l
}

func (l *Line) DebounceTimer() time.Duration {
	return l.debounceTime
}

func (l *Line) Watch(edge Edge, h func(*gpio.Pin)) error {
	l.handler = h
	l.edge = edge
	l.gpioPin.Mode()
	return l.gpioPin.Watch(gpio.Edge(edge), handler)
}

func (l *Line) Unwatch() {
	l.gpioPin.Unwatch()
}

func (l *Line) TestPin(edge Edge) {
}

func (l *Line) Input() {
	l.gpioPin.Input()
}

func (l *Line) PullUp() {
	l.gpioPin.PullUp()
}

func (l *Line) PullDown() {
	l.gpioPin.PullDown()
}

func (l *Line) Pin() int {
	return l.gpioPin.Pin()
}

func (l *Line) Read() bool {
	return bool(l.gpioPin.Read())
}

func handler(pin *gpio.Pin) {
	for _, l := range lines {
		if l.gpioPin == pin {
			if l.debounceTime == 0 {
				l.lastLevel = l.gpioPin.Read()
				l.handler(pin)
				return
			}

			select {
			case <-l.debounceTimer.C:
				l.debounceTimer.Reset(l.debounceTime)
			default:
				return
			}

			go func(p *gpio.Pin, l Line) {
				<-l.debounceTimer.C
				l.debounceTimer.Reset(0)

				switch l.edge {
				case EdgeBoth:
					if p.Read() != l.lastLevel {
						l.lastLevel = l.gpioPin.Read()
						l.handler(pin)
					}
				case EdgeFalling:
					if p.Read() == gpio.Low {
						l.lastLevel = l.gpioPin.Read()
						l.handler(pin)
					}
				case EdgeRising:
					if p.Read() == gpio.High {
						l.lastLevel = l.gpioPin.Read()
						l.handler(pin)
					}
				}
			}(pin, *l)
			return
		}
	}
}
