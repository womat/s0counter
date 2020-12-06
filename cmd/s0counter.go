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

type s0value struct {
	counter     int64     // s0 counter since program start
	lastCounter int64     // s0 counter at the last average calculation
	timeStamp   time.Time // time of last s0 pulse

}

type meter struct {
	//	sync.Mutex
	Gpio          int       // gpio pin
	ScaleFactor   float64   // eg. 1000pulse = 1kWh
	timeStamp     time.Time // time of last throughput calculation
	measuredValue float64   // current measured value, eg kWh,
	throughput    float64   //
	s0            s0value
}

type Save struct {
	MeasuredValue float64   // current measured value, eg kWh,
	TimeStamp     time.Time // time of last s0 pulse
}

type SaveAll map[string]Save

var AllMeters = map[string]meter{}

func main() {
	for meterName, meterConfig := range global.Config.Meter {
		AllMeters[meterName] = meter{
			Gpio:        meterConfig.Gpio,
			ScaleFactor: meterConfig.ScaleFactor,
		}
	}

	err := LoadMeasurements(global.Config.DataFile, AllMeters)
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

	go calcAverage(AllMeters, global.Config.TimerPeriod)

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
	go saveMeasurements(global.Config.DataFile, AllMeters, global.Config.TimerPeriod)
	// wait for a kill signal

	for range time.Tick(10 * time.Second) {
		debug.InfoLog.Println(AllMeters)
	}

	select {}
}

// func handler(pin *gpio.Pin) {
func handlerEmu(pin *rpiemu.PinEmu) {
	p := pin.Pin()

	for name, m := range AllMeters {
		// find the measuring device based on the pin configuration
		if m.Gpio == p {
			// add current counter & set time stamp
			debug.DebugLog.Printf("receive an impule on pin: %v\n", p)
			// m.Lock()
			// defer m.Unlock()

			m.measuredValue += 1 / m.ScaleFactor
			m.s0.counter++
			m.s0.timeStamp = time.Now()
			fmt.Printf("tick port: %v, counter: %v, measuered value: %v,timestamp: %v\n", p, m.s0.counter, m.measuredValue, m.s0.timeStamp)
			AllMeters[name] = m
			return
		}
	}
}

func calcAverage(meters map[string]meter, period time.Duration) {
	for range time.Tick(period) {
		for _, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()

				m.throughput = float64(m.s0.counter-m.s0.lastCounter) / m.ScaleFactor
				m.s0.lastCounter = m.s0.counter
				m.timeStamp = time.Now()
			}()
		}
	}
}

func LoadMeasurements(fileName string, allMeters map[string]meter) (err error) {
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
			meter.measuredValue = loadedMeter.MeasuredValue
			meter.timeStamp = loadedMeter.TimeStamp
			allMeters[name] = meter
		}
	}

	return
}

func saveMeasurements(fileName string, meters map[string]meter, period time.Duration) {
	for range time.Tick(period) {
		var s SaveAll

		for name, m := range meters {
			func() {
				// m.Lock()
				// defer m.Unlock()

				s[name] = Save{
					MeasuredValue: m.measuredValue,
					TimeStamp:     m.timeStamp,
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
