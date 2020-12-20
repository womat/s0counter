package global

import (
	"io"
	"s0counter/pkg/raspberry"
	"sync"
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
//TODO: increase version number to 1.0.1+2020xxyy
const VERSION = "1.0.4+20201220"
const MODULE = "s0counter"

type DebugConf struct {
	File io.WriteCloser
	Flag int
}

type MeterConf struct {
	ScaleFactor float64
	Gpio        int
	BounceTime  time.Duration
}

type WebserverConf struct {
	Port        int             `yaml:"port"`
	Webservices map[string]bool `yaml:"webservices"`
}

type Configuration struct {
	DataCollectionInterval time.Duration
	DataFile               string
	BackupInterval         time.Duration
	Debug                  DebugConf
	Meter                  map[string]MeterConf
	Webserver              WebserverConf
}

type S0 struct {
	Counter     int64     // s0 counter since program start
	LastCounter int64     // s0 counter at the last average calculation
	TimeStamp   time.Time // time of last s0 pulse
}

type Meter struct {
	sync.RWMutex
	LineHandler  *raspberry.Line `yaml:"-"`
	Config       MeterConf
	TimeStamp    time.Time // time of last throughput calculation
	MeterReading float64   // current meter reading (aktueller Zählerstand), eg kWh, l, m³
	FlowPerHour  float64   // mass flow rate per hour (= Flow/time(h)), eg kW, l/h, m³/h
	S0           S0
}

// MetersMap must be a pointer to the Meter type, otherwise RWMutex doesn't work!
type MetersMap = map[string]*Meter

// Config holds the global configuration
var Config Configuration
var AllMeters = MetersMap{}

func init() {
	Config = Configuration{
		Meter:     map[string]MeterConf{},
		Webserver: WebserverConf{Webservices: map[string]bool{}},
	}
}
