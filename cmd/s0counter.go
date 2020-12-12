package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"

	"s0counter/global"
	_ "s0counter/pkg/config"
	"s0counter/pkg/debug"
	"s0counter/pkg/raspberry"
	"s0counter/pkg/tools"
	_ "s0counter/pkg/webservice"
)

type SavedRecord struct {
	MeterReading float64   `yaml:"meterreading"` // current meter reading (aktueller Zählerstand), eg kWh, l, m³
	TimeStamp    time.Time `yaml:"timestamp"`    // time of last s0 pulse
}
type SaveMeters map[string]SavedRecord

func main() {
	debug.SetDebug(global.Config.Debug.File, global.Config.Debug.Flag)

	for meterName, meterConfig := range global.Config.Meter {
		global.AllMeters[meterName] = global.Meter{Config: meterConfig}
	}

	err := loadMeasurements(global.Config.DataFile, global.AllMeters)
	if err != nil {
		debug.FatalLog.Printf("can't open data file: %v\n", err)
		// Exit wit Exit Code 1
		os.Exit(1)
		return
	}

	if err = raspberry.Open(); err != nil {
		debug.FatalLog.Printf("can't open gpio: %v\n", err)
		// Exit wit Exit Code 1
		os.Exit(1)
		return
	}
	defer raspberry.Close()

	for _, meterConfig := range global.Config.Meter {
		pin := raspberry.NewPin(meterConfig.Gpio)
		pin.Input()
		pin.PullUp()
		// call handler when pin changes from low to high.
		pin.Watch(raspberry.EdgeRising, handler)
		defer pin.Unwatch()

		go testPinEmu(pin)
	}

	go calcAverage(global.AllMeters, global.Config.DataCollectionInterval)
	go saveMeasurements(global.Config.DataFile, global.AllMeters, global.Config.BackupInterval)

	// wait for a kill signal
	select {}
}

func calcAverage(meters map[string]global.Meter, period time.Duration) {
	for range time.Tick(period) {
		debug.DebugLog.Println("calc average values")
		for name, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()
				m.FlowPerHour = float64(m.S0.Counter-m.S0.LastCounter) / period.Hours() * m.Config.ScaleFactor
				m.S0.LastCounter = m.S0.Counter
				m.TimeStamp = time.Now()
				meters[name] = m
			}()
		}
	}
}

func loadMeasurements(fileName string, allMeters map[string]global.Meter) (err error) {
	// if file doesn't exists, create an empty file
	if !tools.FileExists(fileName) {
		s := SaveMeters{}

		for name := range allMeters {
			s[name] = SavedRecord{MeterReading: 0, TimeStamp: time.Time{}}
		}

		// marshal the byte slice which contains the yaml file's content into SaveMeters struct
		var data []byte
		data, err = yaml.Marshal(&s)
		if err != nil {
			return
		}

		if err = ioutil.WriteFile(fileName, data, 0600); err != nil {
			return
		}
	}

	// read the yaml file as a byte array.
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	// unmarshal the byte slice which contains the yaml file's content into SaveMeters struct
	s := SaveMeters{}
	if err = yaml.Unmarshal(data, &s); err != nil {
		return
	}

	for name, loadedMeter := range s {
		if meter, ok := allMeters[name]; ok {
			meter.MeterReading = loadedMeter.MeterReading
			meter.TimeStamp = loadedMeter.TimeStamp
			allMeters[name] = meter
		}
	}

	return
}

func saveMeasurements(fileName string, meters map[string]global.Meter, period time.Duration) {
	for range time.Tick(period) {
		debug.DebugLog.Println("save measurements to file")
		s := SaveMeters{}

		for name, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()

				s[name] = SavedRecord{MeterReading: m.MeterReading, TimeStamp: m.TimeStamp}
			}()
		}

		// marshal the byte slice which contains the yaml file's content into SaveMeters struct
		data, err := yaml.Marshal(&s)
		if err != nil {
			debug.ErrorLog.Printf("saveMeasurements marshal: %v\n", err)
			continue
		}

		if err := ioutil.WriteFile(fileName, data, 0600); err != nil {
			debug.ErrorLog.Printf("saveMeasurements write file: %v\n", err)
			continue
		}
	}
}

func increaseImpulse(pin int) {
	for name, m := range global.AllMeters {
		// find the measuring device based on the pin configuration
		if m.Config.Gpio == pin {
			// add current counter & set time stamp
			debug.DebugLog.Printf("receive an impulse on pin: %v\n", pin)
			// m.Lock()
			// defer m.Unlock()
			m.MeterReading += m.Config.ScaleFactor
			m.S0.Counter++
			m.S0.TimeStamp = time.Now()
			global.AllMeters[name] = m
			return
		}
	}
}
