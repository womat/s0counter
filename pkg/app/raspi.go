package app

import (
	"s0counter/pkg/raspberry"
	"time"

	"github.com/womat/debug"
)

// testPinEmu emulate ticks on gpio pin, only for testing in windows mode
func testPinEmu(p raspberry.Pin) {
	for range time.Tick(time.Duration(p.Pin()/2) * time.Second) {
		p.EmuEdge(raspberry.EdgeFalling)
	}
}

func (app *App) handler(p raspberry.Pin) {
	pin := p.Pin()

	for _, m := range app.meters {
		// find the measuring device based on the pin configuration
		if m.Config.Gpio == pin {
			// add current counter & set time stamp
			debug.TraceLog.Printf("receive an impulse on pin: %v", pin)

			t := time.Now()
			m.Lock()
			m.S0.LastTimeStamp = m.S0.TimeStamp
			m.S0.TimeStamp = t
			m.S0.Tick++
			m.Unlock()
			return
		}
	}
}
