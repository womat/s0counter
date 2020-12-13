package raspberry

// Edge represents the change in Line level that triggers an interrupt.
type Edge string

const (
	// EdgeNone indicates no level transitions will trigger an interrupt
	EdgeNone Edge = "none"

	// EdgeRising indicates an interrupt is triggered when the Line transitions from low to high.
	EdgeRising Edge = "rising"

	// EdgeFalling indicates an interrupt is triggered when the Line transitions from high to low.
	EdgeFalling Edge = "falling"

	// EdgeBoth indicates an interrupt is triggered when the Line changes level.
	EdgeBoth Edge = "both"
)

/*
func (pin *Line) Polling(edge Edge, ) {
	lastState := pin.gpioPin.Read()

	for ; true; <-time.After(500 * time.Millisecond) {
		if p := pin.gpioPin.Read(); p != lastState {
			debug.InfoLog.Printf("pin %v is %v\n", pin.Pin(), p)

			switch edge {
			case EdgeBoth:
				debug.InfoLog.Printf("pin %v switch from %v to %v\n", pin.Pin(), lastState, p)
			case EdgeFalling:
				if !p {
					debug.InfoLog.Printf("pin %v switch to (Low) %v\n", pin.Pin(), p)
				}
			case EdgeRising:
				if p {
					debug.InfoLog.Printf("pin %v switch to (High) %v\n", pin.Pin(), p)
				}
			}
			lastState = p
		}
	}
}
*/
