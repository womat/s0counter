package tools

import (
	"testing"
	"time"
)

func TestGetConnectionDeviceIdTimeOut(t *testing.T) {
	testPattern := []struct {
		pattern    string
		connection string
		deviceId   uint8
		timeOut    time.Duration
		maxRetries int
	}{
		{"TCP 192.168.65.197:502 DeviceId:2 Timeout:100 MaxRetries:2", "192.168.65.197:502", 2, 100 * time.Millisecond, 2},
		{"HTTP http://fritz.box Timeout:500", "http://fritz.box", 1, 500 * time.Millisecond, 0},
		{"HTTP https://fritz.box/abc?t=3 Timeout:500 MaxRetries:1", "https://fritz.box/abc?t=3", 1, 500 * time.Millisecond, 1},
	}

	for _, test := range testPattern {
		c, d, to, r := GetConnectionDeviceIdTimeOut(test.pattern)
		if !IsEqual(test.connection, c) || !IsEqual(test.deviceId, d) || !IsEqual(test.timeOut, to) || !IsEqual(test.maxRetries, r) {
			t.Errorf("expected %v %v %v %v, got %v %v %v %v", test.connection, test.deviceId, test.timeOut, test.maxRetries, c, d, to, r)
		}
	}
}

func TestPortSerialTimeOut(t *testing.T) {
	testPattern := []struct {
		pattern  string
		port     string
		baudRate uint
		dataBits uint
		parity   string
		stopBit  uint
		timeOut  time.Duration
	}{
		{"RTU com3,9600,8,O,1 Timeout:5000", "com3", 9600, 8, "O", 1, 5000 * time.Millisecond},
		{"/dev/ttyS0,19200,8,N,2 Timeout:1000", "/dev/ttyS0", 19200, 8, "N", 2, 1000 * time.Millisecond},
	}

	for _, test := range testPattern {
		port, b, d, p, s, to := GetPortSerialTimeOut(test.pattern)
		if !IsEqual(test.port, port) || !IsEqual(test.baudRate, b) || !IsEqual(test.dataBits, d) || !IsEqual(test.parity, p) || !IsEqual(test.stopBit, s) || !IsEqual(test.timeOut, to) {
			t.Errorf("expected %v %v %v %v %v %v, got %v %v %v %v %v %v", test.port, test.baudRate, test.dataBits, test.parity, test.stopBit, test.timeOut, port, b, d, p, s, to)
		}
	}
}
