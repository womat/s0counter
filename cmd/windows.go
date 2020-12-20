// +build windows

package main

import (
	"time"

	"s0counter/pkg/raspberry"
)

func testPinEmu(l *raspberry.Line) {
	for range time.Tick(time.Duration(l.Pin()/2) * time.Second) {
		l.TestPin(raspberry.EdgeFalling)
	}
}

func handler(l *raspberry.Line) {
	increaseImpulse(l.Pin())
}
