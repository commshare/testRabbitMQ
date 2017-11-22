package amqpClient

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var (
	errCantReadConfigFile   = "Unable to read the configuration file : %s "
	errCantDecodeConfigFile = "Unable to decode the configuration file : %s "
)

//newAppConfiguration parses the and convert the yaml file to a configuration object
//P.S. there are not any default configuration to fall back to, so all parameters in the config.yaml file are required

//look a the sample examples/config.yml file for details
func newAppConfiguration(configFile string) (*Config, error) {
	filePath, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf(errCantReadConfigFile, configFile)
	}
	var c Config

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, fmt.Errorf(errCantDecodeConfigFile, configFile)
	}
	return &c, nil
}
