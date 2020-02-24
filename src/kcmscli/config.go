/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * config file
*/
package main

import (
	"os"
	"fmt"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"github.com/kirsle/configdir"
	"github.com/pborman/getopt/v2"
)

// config file
type ConfigFile struct {
	ApiHost string `yaml:"api_host"`
	RefreshToken string `yaml:"refresh_token"`
}

// read & return config file
// if config does not exists returns empty ConfigFile
func readConfigFile() ConfigFile {
	// attempt to read config file
	path := filepath.Join(configdir.LocalConfig("kcmscli"), "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ConfigFile{}
	}
	configYAML := readBytes(path)
	
	// parse config file yaml
	var config ConfigFile
	yaml.Unmarshal(configYAML, &config)

	return config
}

// commit changes to config file to disk
func (config ConfigFile) commit() {
	// create config dir if does already exists
	dir := configdir.LocalConfig("kcmscli")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	path := filepath.Join(dir, "config.yaml")

	// convert config to yaml
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		die(err.Error())
	}
	
	writeBytes(configYAML, path)
}

// config subcommand
// globalConfig - global program config 
// args - arguments parsed to subcommand
func configCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli config [options] <item> [value]
kcmscli config  - change config item value to given value

if [value] is obmitted, would get value of item instead

CONFIG
ApiHost - set the k8s-cms master api host endpoint used to communciate with k8s-cms

OPTIONS
`
	// parse & evaluate options
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp , "help", 'h', "show usage info")
	optSet.FlagLong(&globalConfig.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// parse positional args
	args = optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <item> [value]")
	}
	configItem, configValue := args[0], ""
	if len(args) > 1 {
		configValue = args[1]
	}

	// evaluate config cmd - get/set
	configFile := readConfigFile()
	switch configItem {
	case "ApiHost":
		if configValue == "" { 
			// get value
			fmt.Println(configFile.ApiHost)
		} else {
			configFile.ApiHost = configValue
		}
	default:
		fmt.Printf("Unknown config item: %s\n", configItem)
		os.Exit(1);
	}
	
	configFile.commit()
}
