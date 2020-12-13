package main

//TODO: move SavedRecord to package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
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
		global.AllMeters[meterName] = &global.Meter{Config: meterConfig}
	}

	err := loadMeasurements(global.Config.DataFile, global.AllMeters)
	if err != nil {
		debug.FatalLog.Printf("can't open data file: %v\n", err)
		os.Exit(1)
		return
	}

	if err = raspberry.Open(); err != nil {
		debug.FatalLog.Printf("can't open gpio: %v\n", err)
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

	go calcFlowPerHour(global.AllMeters, global.Config.DataCollectionInterval)
	go backupMeasurements(global.Config.DataFile, global.AllMeters, global.Config.BackupInterval)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// wait for am os.Interrupt signal (CTRL C)
	sig := <-c
	debug.InfoLog.Printf("Got %s signal. Aborting...\n", sig)
	_ = saveMeasurements(global.Config.DataFile, global.AllMeters)
	os.Exit(1)
}

func calcFlowPerHour(meters global.MetersMap, period time.Duration) {
	for range time.Tick(period) {
		debug.DebugLog.Println("calc average values")

		for _, m := range meters {
			func() {
				m.Lock()
				defer m.Unlock()
				m.FlowPerHour = float64(m.S0.Counter-m.S0.LastCounter) / period.Hours() * m.Config.ScaleFactor
				m.S0.LastCounter = m.S0.Counter
				m.TimeStamp = time.Now()
			}()
		}
	}
}

func loadMeasurements(fileName string, allMeters global.MetersMap) (err error) {
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
			func() {
				meter.Lock()
				defer meter.Unlock()
				meter.MeterReading = loadedMeter.MeterReading
				meter.TimeStamp = loadedMeter.TimeStamp
			}()
		}
	}

	return
}

func backupMeasurements(fileName string, meters global.MetersMap, period time.Duration) {
	for range time.Tick(period) {
		_ = saveMeasurements(fileName, meters)
	}
}

func saveMeasurements(fileName string, meters global.MetersMap) error {
	debug.DebugLog.Println("saveMeasurements measurements to file")

	s := SaveMeters{}

	for name, m := range meters {
		func() {
			m.RLock()
			defer m.RUnlock()
			s[name] = SavedRecord{MeterReading: m.MeterReading, TimeStamp: m.TimeStamp}
		}()
	}

	// marshal the byte slice which contains the yaml file's content into SaveMeters struct
	data, err := yaml.Marshal(&s)
	if err != nil {
		debug.ErrorLog.Printf("backupMeasurements marshal: %v\n", err)
		return err
	}

	if err := ioutil.WriteFile(fileName, data, 0600); err != nil {
		debug.ErrorLog.Printf("backupMeasurements write file: %v\n", err)
		return err
	}

	return nil
}

func increaseImpulse(pin int) {
	for _, m := range global.AllMeters {
		// find the measuring device based on the pin configuration
		if m.Config.Gpio == pin {
			// add current counter & set time stamp
			debug.DebugLog.Printf("receive an impulse on pin: %v\n", pin)

			func() {
				m.Lock()
				defer m.Unlock()
				m.MeterReading += m.Config.ScaleFactor
				m.S0.Counter++
				m.S0.TimeStamp = time.Now()
			}()

			return
		}
	}
}
