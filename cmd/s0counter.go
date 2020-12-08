package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"s0counter/pkg/rpi"
	"s0counter/pkg/tools"
	"time"

	"s0counter/global"
	_ "s0counter/pkg/config"
	"s0counter/pkg/debug"
	_ "s0counter/pkg/webservice"
)

type Save struct {
	MeasuredValue float64   // current measured value, eg kWh,
	TimeStamp     time.Time // time of last s0 pulse
}

type SaveAll map[string]struct {
	MeasuredValue float64   // current measured value, eg kWh,
	TimeStamp     time.Time // time of last s0 pulse
}

func main() {
	debug.SetDebug(global.Config.Debug.File, global.Config.Debug.Flag)

	for meterName, meterConfig := range global.Config.Meter {
		global.AllMeters[meterName] = global.Meter{Config: meterConfig}
	}

	err := LoadMeasurements(global.Config.DataFile, global.AllMeters)
	if err != nil {
		debug.FatalLog.Printf("can't open data file: %v\n", err)
		// Exit wit Exit Code 1
		os.Exit(1)
		return
	}

	rb, err := rpi.Open()

	if err != nil {
		debug.FatalLog.Printf("can't open gpio: %v\n", err)
		// Exit wit Exit Code 1
		os.Exit(1)
		return
	}
	defer rb.Close()

	for _, meterConfig := range global.Config.Meter {
		pin := rb.NewPin(meterConfig.Gpio)
		pin.Input()
		pin.PullUp()
		// call handler when pin changes from low to high.
		pin.Watch(rpi.EdgeRising, handler)
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
				m.Throughput = float64(m.S0.Counter-m.S0.LastCounter) / period.Hours() * m.Config.ScaleFactor
				debug.DebugLog.Printf("m.Throughput %v, m.S0.Counter %v,m.S0.LastCounter %v, m.Config.ScaleFactor %v, period %v, period.Hours() %v\n", m.Throughput, m.S0.Counter, m.S0.LastCounter, m.Config.ScaleFactor, period, period.Hours())
				m.S0.LastCounter = m.S0.Counter
				m.TimeStamp = time.Now()
				meters[name] = m
			}()
		}
	}
}

func LoadMeasurements(fileName string, allMeters map[string]global.Meter) (err error) {
	// if file doesn't exists, create an empty file
	if !tools.FileExists(fileName) {
		s := SaveAll{}

		for name := range allMeters {
			s[name] = Save{}
		}

		// marshal the byte slice which contains the yaml file's content into SaveAll struct
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

	// unmarshal the byte slice which contains the yaml file's content into SaveAll struct
	s := SaveAll{}
	if err = yaml.Unmarshal(data, &s); err != nil {
		return
	}

	for name, loadedMeter := range s {
		if meter, ok := allMeters[name]; ok {
			meter.MeasuredValue = loadedMeter.MeasuredValue
			meter.TimeStamp = loadedMeter.TimeStamp
			allMeters[name] = meter
		}
	}

	return
}

func saveMeasurements(fileName string, meters map[string]global.Meter, period time.Duration) {
	for range time.Tick(period) {
		debug.DebugLog.Println("save measurements to file")
		s := SaveAll{}

		for name, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()

				s[name] = Save{
					MeasuredValue: m.MeasuredValue,
					TimeStamp:     m.TimeStamp,
				}
			}()
		}

		// marshal the byte slice which contains the yaml file's content into SaveAll struct
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
