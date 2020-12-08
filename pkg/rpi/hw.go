// +build !windows

package rpi

import (
	"github.com/warthog618/gpio"
)

func Open() (*Rpi, error) {
	return &Rpi{}, gpio.Open()
}

type Rpi struct {
}

type PinHw struct {
	gpioPin *gpio.Pin
}

func (t *Rpi) Close() {
	return
}

func (t *Rpi) NewPin(p int) *PinHw {
	return &(PinHw{gpioPin: gpio.NewPin(p)})
}

func (p *PinHw) Watch(edge Edge, handler interface{}) {
	if tmp, ok := handler.(func(*gpio.Pin)); ok {
		p.gpioPin.Watch(gpio.Edge(edge), tmp)
	}
}

func (p *PinHw) Unwatch() {
	p.gpioPin.Unwatch()
}

func (p *PinHw) TestPin(edge Edge) {
}

func (p *PinHw) Input() {
	p.gpioPin.Input()
}

func (p *PinHw) PullUp() {
	p.gpioPin.PullUp()
}

func (p *PinHw) PullDown() {
	p.gpioPin.PullDown()
}

func (p *PinHw) Pin() int {
	return p.gpioPin.Pin()
}
