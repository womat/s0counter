// +build !windows

package rpihw

import (
	"github.com/warthog618/gpio"
	"s0counter/pkg/raspberry"
)

func Open() (*Rpi, error) {
	return &Rpi{}, gpio.Open()
}

type Rpi struct {
}

type Pin struct {
	gpioPin *gpio.Pin
}

func (t *Rpi) Close() {
	return
}

func (t *Rpi) NewPin(p int) *Pin {
	return &(Pin{gpioPin: gpio.NewPin(p)})
}

func (p *Pin) Watch(edge raspberry.Edge, handler interface{}) {
	if tmp, ok := handler.(func(*Pin)); ok {
		p.gpioPin.Watch(edge, tmp)
	}
}

func (p *Pin) Unwatch() {
	p.gpioPin.Unwatch()
}

func (p *Pin) TestPin(edge raspberry.Edge) {
}

func (p *Pin) Input() {
	p.gpioPin.Input()
}

func (p *Pin) PullUp() {
	p.gpioPin.PullUp()
}

func (p *Pin) PullDown() {
	p.gpioPin.PullDown()
}

func (p *Pin) Pin() int {
	return p.gpioPin.Pin()
}
