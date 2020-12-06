package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"s0counter/global"
	"s0counter/pkg/tools"
)

func init() {
	type yamlStruct struct {
		TimePeriod int
		DataFile   string
		Debug      struct {
			File string
			Flag string
		}
		Meter     map[string]global.Meter
		Webserver global.WebserverConf
	}

	var configFile yamlStruct

	flag.Bool("version", false, "print version and exit")
	flag.String("debug.file", "stderr", "log file eg. /tmp/emu.log")
	flag.String("debug.flag", "", "enable debug information (standard | trace | debug)")
	flag.String("config", "", "Config File eg. /opt/womat/s0counter.yaml")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("version") {
		fmt.Printf("Version: %v\n", global.VERSION)
		os.Exit(0)
	}

	if f := viper.GetString("config"); f != "" {
		viper.SetConfigFile(f)
	} else {
		viper.SetConfigName("s0counter")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/opt/womat/")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configFile)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	getDebugFlag := func(flag string) int {
		switch flag {
		case "trace":
			return Full
		case "debug":
			return Warning | Info | Error | Fatal | Debug
		case "standard":
			return Standard
		}
		return 0
	}

	global.Config.Debug.Flag = getDebugFlag(configFile.Debug.Flag)
	switch file := configFile.Debug.File; file {
	case "stderr":
		global.Config.Debug.File = os.Stderr
	case "stdout":
		global.Config.Debug.File = os.Stdout
	default:
		if !tools.FileExists(file) {
			_ = tools.CreateFile(file)
		}
		if global.Config.Debug.File, err = os.Open(file); err != nil {
			fatalLog.Println(err)
			os.Exit(0)
		}
	}

	global.Config.Meter = configFile.Meter
	global.Config.DataFile = configFile.DataFile
	global.Config.TimerPeriod = 5 * time.Second
	if configFile.TimePeriod > 0 {
		global.Config.TimerPeriod = time.Duration(configFile.TimePeriod) * time.Second
	}

	global.Config.Webserver = configFile.Webserver
}
