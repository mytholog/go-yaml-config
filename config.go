package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
)

//LoadConfigFromFile loading config from yaml file
func LoadConfigFromFile(config interface{}, configFile string, defaultValue interface{}) string {
	//log.Debugf("Reading configuration from '%s'", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		//log.Warn("Configuration not found")
		if defaultValue != nil {
			//log.Warn("Default value is defined. Using it.")
			copier.Copy(config, defaultValue)
			return ""
		}
		panic(err)
	}

	if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), config); err != nil {
		//log.Warn("Configuration incorrect ")
		if defaultValue != nil {
			//log.Warn("Default value is defined. Use it.")
			copier.Copy(config, defaultValue)
			return ""
		}
		panic(err)
	}

	customConfigFile := filepath.Join(
		filepath.Dir(configFile),
		strings.TrimSuffix(filepath.Base(configFile), filepath.Ext(configFile))+".custom"+filepath.Ext(configFile),
	)
	//log.Debugf("Try to read custom configuration from '%s'...", customConfigFile)
	data, err = ioutil.ReadFile(customConfigFile)
	if err == nil {
		//log.Debugf("Reading custom configuration from '%s'", customConfigFile)
		if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), config); err != nil {
			panic(err)
		}
		//log.Debug("Config loaded successfully with custom config file")
		return customConfigFile
	}

	//log.Debug("Config loaded successfully")
	return configFile
}

//LoadConfig from command line argument
func LoadConfig(config interface{}, defaultFilename string, defaultValue interface{}) string {
	var configFile string
	flag.StringVar(&configFile, "c", defaultFilename, "Config file")
	flag.StringVar(&configFile, "config", defaultFilename, "Config file")
	flag.Parse()
	return LoadConfigFromFile(config, configFile, defaultValue)
}
