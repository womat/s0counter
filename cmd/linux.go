// +build !windows

package main

import (
	"github.com/warthog618/gpio"
	"s0counter/global"
	"s0counter/pkg/debug"
	"s0counter/pkg/rpi"
	"time"
)

func testPinEmu(p *rpi.PinHw) {
}

func handler(pin *gpio.Pin) {
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
