/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
*/
package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
)


// global program config - common config for all subcommands
type GlobalConfig struct {
	isVerbose bool
	shouldHelp bool
}


func main() {
	var usageInfo string = `Usage: kcmscli [options] <subcommand ...> 
kcmscli - Command line tool for controlling k8s-cms via k8s-cms master

SUBCOMMANDS
* use <subcommand> -h to show usage info for each subcommand
config - change configuration values
auth - authenticate, deauthenticate with k8s-cms master.
user - import, export users
contest - import, list, delete contests

OPTIONS
`
	// parse & evaluate program opts
	var config GlobalConfig
	optSet := getopt.New()
	optSet.FlagLong(&config.shouldHelp , "help", 'h', "show usage info")
	optSet.FlagLong(&config.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.Parse(os.Args)

	if config.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// parse subcommand
	args := optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <subcommand>")
	}
	subCmd := args[0]

	// delegate to various subcommands
	switch subCmd {
	case "config":
		configCmd(&config, args)
	case "auth" :
		authCmd(&config, args)
	case "user":
		userCmd(&config, args)
	case "contest":
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCmd)
		os.Exit(1);
	}
}
