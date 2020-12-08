// +build windows

package main

import (
	"time"

	"s0counter/pkg/raspberry"
)

func testPinEmu(p *raspberry.P) {
	for range time.Tick(time.Duration(p.Pin()/2) * time.Second) {
		p.TestPin(raspberry.EdgeRising)
	}
}

func handler(pin *raspberry.P) {
	increaseImpulse(pin.Pin())
}
