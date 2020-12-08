// +build windows

package main

import (
	"s0counter/global"
	"s0counter/pkg/debug"
	"s0counter/pkg/rpi"
	"time"
)

func testPinEmu(p *rpi.PinEmu) {
	for range time.Tick(time.Duration(p.Pin()/2) * time.Second) {
		p.TestPin(rpi.EdgeRising)
	}
}

func handler(pin *rpi.PinEmu) {
	p := pin.Pin()

	for name, m := range global.AllMeters {
		// find the measuring device based on the pin configuration
		if m.Config.Gpio == p {
			// add current counter & set time stamp
			debug.DebugLog.Printf("receive an impulse on pin: %v\n", p)
			// m.Lock()
			// defer m.Unlock()
			m.MeasuredValue += 1 / m.Config.ScaleFactor
			m.S0.Counter++
			m.S0.TimeStamp = time.Now()
			global.AllMeters[name] = m
			return
		}
	}
}
