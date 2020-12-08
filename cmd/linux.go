// +build !windows

package main

import (
	"github.com/warthog618/gpio"

	"s0counter/pkg/raspberry"
)

func testPinEmu(p *raspberry.P) {
}

func handler(pin *gpio.Pin) {
	increaseImpulse(pin.Pin())
}
