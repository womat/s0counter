package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"s0counter/pkg/rpiemu"
	"s0counter/pkg/tools"
	"time"

	"s0counter/global"
	_ "s0counter/pkg/config"
	"s0counter/pkg/debug"
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

	rb, err := rpiemu.Open()

	if err != nil {
		debug.FatalLog.Printf("can't open gpio: %v\n", err)
		// Exit wit Exit Code 1
		os.Exit(1)
		return
	}
	defer rb.Close()

	go calcAverage(global.AllMeters, global.Config.DataCollectionInterval)

	for _, meterConfig := range global.Config.Meter {
		pin := rb.NewPin(meterConfig.Gpio)
		pin.Input()
		pin.PullUp()
		// call handler when pin changes from low to high.
		// TODO: determine Hw and emu
		pin.Watch(rpiemu.EdgeRising, handlerEmu)
		defer pin.Unwatch()

		go func(p *rpiemu.PinEmu) {
			for range time.Tick(time.Duration(pin.Pin()/2) * time.Second) {
				p.TestPin(rpiemu.EdgeRising)
			}
		}(pin)
	}

	time.Sleep(time.Second)
	go saveMeasurements(global.Config.DataFile, global.AllMeters, global.Config.BackupInterval)
	// wait for a kill signal

	for range time.Tick(10 * time.Second) {
		debug.InfoLog.Println(global.AllMeters)
	}

	select {}
}

// func handler(pin *gpio.Pin) {
func handlerEmu(pin *rpiemu.PinEmu) {
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
			fmt.Printf("tick port: %v, counter: %v, measuered value: %v,timestamp: %v\n", p, m.S0.Counter, m.MeasuredValue, m.S0.TimeStamp)
			global.AllMeters[name] = m
			return
		}
	}
}

func calcAverage(meters map[string]global.Meter, period time.Duration) {
	for range time.Tick(period) {
		for _, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()

				m.Throughput = float64(m.S0.Counter-m.S0.LastCounter) / m.Config.ScaleFactor
				m.S0.LastCounter = m.S0.Counter
				m.TimeStamp = time.Now()
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
		var s SaveAll

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
