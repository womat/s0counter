package rpiemu

// Edge represents the change in PinEmu level that triggers an interrupt.
type Edge string

const (
	// EdgeNone indicates no level transitions will trigger an interrupt
	EdgeNone Edge = "none"

	// EdgeRising indicates an interrupt is triggered when the PinEmu transitions from low to high.
	EdgeRising Edge = "rising"

	// EdgeFalling indicates an interrupt is triggered when the PinEmu transitions from high to low.
	EdgeFalling Edge = "falling"

	// EdgeBoth indicates an interrupt is triggered when the PinEmu changes level.
	EdgeBoth Edge = "both"
)

// Framer is the interface that handle gpio
type Gpio interface {
	Close()
	//	NewPin(int) *PinEmu
}

// Pins is the interface that handles Pins.
type Pins interface {
	Watch(Edge, interface{})
	Unwatch()
	TestPin(Edge)
	Input()
	PullUp()
	PullDown()
	Pin() int
}

func Close(t Gpio) {
	t.Close()
}

func Watch(p Pins, edge Edge, handler interface{}) {
	p.Watch(edge, handler)
}

func Unwatch(p Pins) {
	p.Unwatch()
}

func Input(p Pins) {
	p.Input()
}

func PullUp(p Pins) {
	p.PullUp()
}

func PullDown(p Pins) {
	p.PullDown()
}

func Pin(p Pins) int {
	return p.Pin()
}
