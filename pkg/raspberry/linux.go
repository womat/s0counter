// +build !windows

package raspberry

import (
	"github.com/warthog618/gpio"
	"s0counter/pkg/debug"
	"time"
)

type P struct {
	gpioPin *gpio.Pin
}
type C struct {
	gpioPin *gpio.Pin
}

func Open() error {
	c, err := gpiod.NewChip("gpiochip0")

	return gpio.Open()
}

func Close() error {
	debug.InfoLog.Printf("gpio Close\n")
	return gpio.Close()
}

func NewPin(p int) *P {
	return &(P{gpioPin: gpio.NewPin(p)})
}

func (p *P) Watch(edge Edge, handler func(*gpio.Pin)) error {
	//return p.gpioPin.Watch(gpio.Edge(edge), handler)
	debug.InfoLog.Printf("watch pin %v\n", p.Pin())

	go func(edge gpio.Edge, pin *gpio.Pin) {
		lastState := pin.Read()

		for ; true; <-time.After(500 * time.Millisecond) {
			if p := pin.Read(); p != lastState {
				debug.InfoLog.Printf("pin %v is %v\n", pin.Pin(), p)

				switch edge {
				case gpio.EdgeBoth:
					debug.InfoLog.Printf("pin %v switch from %v to %v\n", pin.Pin(), lastState, p)
				case gpio.EdgeFalling:
					if !p {
						debug.InfoLog.Printf("pin %v switch to (Low) %v\n", pin.Pin(), p)
					}
				case gpio.EdgeRising:
					if p {
						debug.InfoLog.Printf("pin %v switch to (High) %v\n", pin.Pin(), p)
					}
				}
				lastState = p
			}
		}
	}(gpio.Edge(edge), p.gpioPin)

	p.gpioPin.Unwatch()

	err := p.gpioPin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
		debug.InfoLog.Printf("watcher: pin %v is %v\n", pin.Pin(), pin.Read())
	})

	debug.InfoLog.Printf("returncode watcher pin %v: %v\n", p.Pin(), err)

	return err

}

func (p *P) Unwatch() {
	debug.InfoLog.Printf("unwatch pin %v\n", p.Pin())
	p.gpioPin.Unwatch()
}

func (p *P) TestPin(edge Edge) {
}

func (p *P) Input() {
	debug.InfoLog.Printf("set pin %v to input\n", p.Pin())
	p.gpioPin.Input()
}

func (p *P) PullUp() {
	debug.InfoLog.Printf("pullup pin %v\n", p.Pin())
	p.gpioPin.PullUp()
}

func (p *P) PullDown() {
	debug.InfoLog.Printf("pulldown pin %v\n", p.Pin())
	p.gpioPin.PullDown()
}

func (p *P) Pin() int {
	return p.gpioPin.Pin()
}
