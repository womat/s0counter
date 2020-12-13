// +build !windows

package raspberry

import (
	"github.com/warthog618/gpio"
)

type Line struct {
	gpioPin *gpio.Pin
}
type Chip struct {
}

func Open() (*Chip, error) {
	return &Chip{}, gpio.Open()
}

func (c *Chip) Close() (err error) {
	return gpio.Close()
}

func (c *Chip) NewPin(p int) *Line {
	return &(Line{gpio.NewPin(p)})
}

func (l *Line) Watch(edge Edge, handler func(*gpio.Pin)) error {
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
