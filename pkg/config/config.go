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
)

// defaultInterval defines the default of DataCollectionInterval and BackupInterval
const defaultInterval = 60 * time.Second

func init() {
	type yamlStruct struct {
		DataCollectionInterval int
		DataFile               string
		BackupInterval         int

		Debug struct {
			File string
			Flag string
		}
		Meter map[string]struct {
			ScaleFactor float64
			Gpio        int
		}
		Webserver struct {
			Port        int
			Webservices map[string]bool
		}
	}

	var configFile yamlStruct

	flag.Bool("version", false, "print version and exit")
	flag.String("debug.file", "stderr", "log file eg. /opt/womat/log/"+global.MODULE+".log")
	flag.String("debug.flag", "", "enable debug information (standard | trace | debug)")
	flag.String("config", "/opt/womat/config/"+global.MODULE+".yaml", "Config File eg. /opt/womat/config/"+global.MODULE+"+.yaml")

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
		viper.SetConfigName(global.MODULE)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/opt/womat/config")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configFile)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// defines Debug section of global.Config
	global.Config.Debug.Flag = getDebugFlag(configFile.Debug.Flag)
	switch file := configFile.Debug.File; file {
	case "stderr":
		global.Config.Debug.File = os.Stderr
	case "stdout":
		global.Config.Debug.File = os.Stdout
	default:
		if global.Config.Debug.File, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			fatalLog.Println(err)
			os.Exit(1)
		}
	}

	// defines dataFile, backupInterval and DataCollectionInterval of global.Config
	global.Config.DataFile = configFile.DataFile

	if configFile.BackupInterval > 0 {
		global.Config.BackupInterval = time.Duration(configFile.BackupInterval) * time.Second
	} else {
		global.Config.BackupInterval = defaultInterval
	}
	if configFile.DataCollectionInterval > 0 {
		global.Config.DataCollectionInterval = time.Duration(configFile.DataCollectionInterval) * time.Second
	} else {
		global.Config.DataCollectionInterval = defaultInterval
	}

	// defines the Meter section of global.Config
	for meterName, meterConfig := range configFile.Meter {
		global.Config.Meter[meterName] = meterConfig
	}

	// defines the WebServer section of global.Config
	global.Config.Webserver = configFile.Webserver
}

func getDebugFlag(flag string) int {
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
