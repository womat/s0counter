// +build !windows

package raspberry

import (
	"github.com/warthog618/gpio"
)

type P struct {
	gpioPin *gpio.Pin
}

func Open() error {
	return gpio.Open()
}

func Close() error {
	return gpio.Close()
}

func NewPin(p int) *P {
	return &(P{gpioPin: gpio.NewPin(p)})
}

func (p *P) Watch(edge Edge, handler func(*gpio.Pin)) error {
	return p.gpioPin.Watch(gpio.Edge(edge), handler)
}

func (p *P) Unwatch() {
	p.gpioPin.Unwatch()
}

func (p *P) TestPin(edge Edge) {
}

func (p *P) Input() {
	p.gpioPin.Input()
}

func (p *P) PullUp() {
	p.gpioPin.PullUp()
}

func (p *P) PullDown() {
	p.gpioPin.PullDown()
}

func (p *P) Pin() int {
	return p.gpioPin.Pin()
}
