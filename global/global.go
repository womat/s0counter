package global

import (
	"io"
	"time"
)

// VERSION holds the version information with the following logic in mind
//  1 ... fixed
//  0 ... year 2020, 1->year 2021, etc.
//  7 ... month of year (7=July)
//  the date format after the + is always the first of the month
//
// VERSION differs from semantic versioning as described in https://semver.org/
// but we keep the correct syntax.
//TODO: Versionsnummer auf 1.0.1+2020xxyy anpassen
const VERSION = "1.0.0+20201123"

type DebugConf struct {
	File io.WriteCloser
	Flag int
}

type Meter struct {
	ScaleFactor int
	Gpio        int
}

type WebserverConf struct {
	Port        int
	Webservices map[string]bool
}

type Configuration struct {
	TimerPeriod time.Duration
	DataFile    string
	Debug       DebugConf
	Meter       map[string]Meter
	Webserver   WebserverConf
}

// Config holds the global configuration
var Config Configuration

func init() {
	Config = Configuration{
		Meter:     map[string]Meter{},
		Webserver: WebserverConf{Webservices: map[string]bool{}},
	}
}
