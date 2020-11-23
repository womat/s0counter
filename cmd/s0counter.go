package main

import (
	"fmt"
	"github.com/warthog618/gpio"
	"s0counter/global"
	_ "s0counter/pkg/config"
	"time"
)

type value struct {
	counter     int64
	lastCounter int64
	throughput  float64
}

type meter struct {
	Gpio        int
	ScaleFactor int
	lastTime    time.Time
	time        time.Time
	measurand   value
}

func main() {
	AllMeters := map[string]meter{}
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	for meterName, meterConfig := range global.Config.Meter {
		pin := gpio.NewPin(meterConfig.Gpio)
		pin.Input()
		pin.PullUp()
		pin.Watch(gpio.EdgeRising, handler) // Call handler when pin changes from Low to High.
		defer pin.Unwatch()

		AllMeters[meterName] = meter{
			Gpio:        meterConfig.Gpio,
			ScaleFactor: meterConfig.ScaleFactor,
		}
	}

}

func handler(pin *gpio.Pin) {
	fmt.Printf("Pin 4 is %v", pin.Read())
}
