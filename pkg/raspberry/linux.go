// +build !windows

package raspberry

import (
	"fmt"
	"github.com/warthog618/gpio"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"s0counter/pkg/debug"
	"syscall"
	"time"
)

type P struct {
	gpioPin *gpio.Pin
}
type Chip struct {
	GpioChip *gpiod.Chip
}

type Line struct {
	GpioLine *gpiod.Line
}

func Open() (*Chip, error) {
	debug.InfoLog.Printf("open gpio\n")
	c, err := gpiod.NewChip("gpiochip0")
	debug.InfoLog.Printf("returncode open gpio: %v\n", err)

	return &(Chip{GpioChip: c}), err
	//return gpio.Open()
}

func (c *Chip) Close() error {
	debug.InfoLog.Printf("gpio Close\n")
	err := c.GpioChip.Close()
	debug.InfoLog.Printf("returncode close: %v\n", err)

	return err
	//return gpio.Close()
}

func (c *Chip) Watch(p int) (*Line, error) {
	debug.InfoLog.Printf("watch pin %v\n", p)
	i, _ := rpi.Pin("GPIO17")
	debug.InfoLog.Printf("watch pin %v\n", i)
	debug.InfoLog.Printf("watch pin %v\n", rpi.GPIO17)

	l, err := c.GpioChip.RequestLine(p,
		gpiod.WithFallingEdge,
		gpiod.WithEventHandler(eventHandler))

	if err != nil {
		fmt.Printf("RequestLine returned error: %s\n", err)
		if err == syscall.Errno(22) {
			fmt.Println("Note that the WithPullUp option requires kernel V5.5 or later - check your kernel version.")
		}
	}
	debug.InfoLog.Printf("returncode watcher pin %v: %v\n", p, err)
	x, _ := l.Info()
	debug.InfoLog.Printf("info of pin %v: %v\n", p, x)

	return &(Line{GpioLine: l}), err
}

func eventHandler(evt gpiod.LineEvent) {
	t := time.Now()
	edge := "rising"
	if evt.Type == gpiod.LineEventFallingEdge {
		edge = "falling"
	}
	fmt.Printf("event:%3d %-7s %s (%s)\n",
		evt.Offset,
		edge,
		t.Format(time.RFC3339Nano),
		evt.Timestamp)
}

func (l *Line) Close() error {
	debug.InfoLog.Printf("gpio Close\n")
	err := l.GpioLine.Close()
	debug.InfoLog.Printf("returncode close: %v\n", err)

	return err
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
