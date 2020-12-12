package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"

	"s0counter/global"
)

// defaultInterval defines the default of dataCollectionInterval and backupInterval (in seconds)
const defaultInterval = 60

type yamlDebug struct {
	File string `yaml:"file"`
	Flag string `yaml:"flag"`
}

type yamlStruct struct {
	DataCollectionInterval int                         `yaml:"datacollectioninterval"`
	DataFile               string                      `yaml:"datafile"`
	BackupInterval         int                         `yaml:"backupinterval"`
	Debug                  yamlDebug                   `yaml:"debug"`
	Meter                  map[string]global.MeterConf `yaml:"meter"`
	Webserver              global.WebserverConf        `yaml:"webserver"`
}

func init() {
	var err error
	var flags = flagm{
		"version":   {flagType: flagBool, usage: "print version and exit", defaultValue: false},
		"debugFile": {flagType: flagString, usage: "log file eg. /opt/womat/log/" + global.MODULE + `.log (default "stderr")`, defaultValue: ""},
		"debugFlag": {flagType: flagString, usage: `"enable debug information (standard | trace | debug) (default "standard")`, defaultValue: ""},
		"config":    {flagType: flagString, usage: "Config File", defaultValue: "/opt/womat/config/" + global.MODULE + ".yaml"},
	}
	var configFile = yamlStruct{
		DataCollectionInterval: defaultInterval,
		DataFile:               "/opt/womat/data/measurement.yaml",
		BackupInterval:         defaultInterval,
		Debug:                  yamlDebug{File: "stderr", Flag: "standard"},
		Meter:                  map[string]global.MeterConf{},
		Webserver:              global.WebserverConf{Port: 4000, Webservices: map[string]bool{"version": false, "currentdata": false}},
	}

	parse(flags)

	if flags.bool("version") {
		fmt.Printf("Version: %v\n", global.VERSION)
		os.Exit(0)
	}

	if err := readConfigFile(flags.string("config"), &configFile); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	if global.Config.Debug, err = getDebugConfig(flags, configFile.Debug); err != nil {
		log.Fatalf("unable to open debug file, %v", err)
	}
	global.Config.DataFile = configFile.DataFile
	global.Config.BackupInterval = time.Duration(configFile.BackupInterval) * time.Second
	global.Config.DataCollectionInterval = time.Duration(configFile.DataCollectionInterval) * time.Second
	global.Config.Webserver = configFile.Webserver
	global.Config.Meter = configFile.Meter
}

func getDebugConfig(flags flagm, d yamlDebug) (c global.DebugConf, err error) {
	var file, flag string

	if s := flags.string("debugFile"); s != "" {
		file = s
	} else {
		file = d.File
	}

	if s := flags.string("debugFlag"); s != "" {
		flag = s
	} else {
		flag = d.Flag
	}

	// defines Debug section of global.Config
	switch flag {
	case "trace":
		c.Flag = Full
	case "debug":
		c.Flag = Warning | Info | Error | Fatal | Debug
	case "standard":
		c.Flag = Standard
	}

	switch file {
	case "stderr":
		c.File = os.Stderr
	case "stdout":
		c.File = os.Stdout
	default:
		if c.File, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			return
		}
	}

	return
}

func readConfigFile(fileName string, c *yamlStruct) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(c); err != nil {
		return err
	}

	return nil
}
