/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * auth subcommand
 */
package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/pborman/getopt/v2"
	"net/http"
	"os"
)

type LoginCreds struct {
	UserName string `json:"username"`
	Passwd   string `json:"password"`
}

// authentication subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func authCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli auth [options] <subcommand ...>
kcmscli auth - authenticate, deauthenticate with k8s-cms API

SUBCOMMANDS
login  - authenticate with k8s-cms API with admin credentials
logout - ends the authenticated session with k8s-cms API.
check - check if the CLI is currently authenticated with the server

OPTIONS
`
	// parse & evaluate options
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.FlagLong(&globalConfig.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// parse subcommand
	if len(args) < 1 {
		die("Missing positional arguments: <subcommand>")
	}
	args = optSet.Args()
	subCmd, restArgs := args[0], args[1:]

	// delegate to various subcommands
	switch subCmd {
	case "login":
		authLoginCmd(globalConfig, restArgs)
	case "logout":
		authLogoutCmd(globalConfig, restArgs)
	case "check":
		authCheckCmd(globalConfig, restArgs)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

// authentication - login subcommmand
// config - global program config
// args - arguments parsed to subcommand
func authLoginCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli auth login [options] <username>
kcmscli auth login - authenticate as the admin with the given username

OPTIONS
`
	// parse & evaluate options
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.FlagLong(&globalConfig.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// parse login credentials
	if len(args) < 1 {
		die("Missing positional arguments: <username>")
	}
	username := args[0]

	fmt.Printf("%s password: ", username)
	passwd, err := gopass.GetPasswd()
	if err != nil {
		die(err.Error())
	}
	loginCreds := LoginCreds{
		UserName: username,
		Passwd:   string(passwd),
	}

	// perform login API call
	configFile := readConfigFile()
	api := makeAPI(configFile, globalConfig)
	statusCode, response := api.callJSON("POST", "auth/login", loginCreds)

	// parse results of API call
	switch statusCode {
	case http.StatusOK:
		fmt.Printf("Logged in as %s\n", username)
	case http.StatusUnauthorized:
		die("Login failed: Check your login credentials")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", statusCode))
	}

	// save response refresh token for later use
	configFile.RefreshToken = response["refreshToken"].(string)
	configFile.commit()
}

// authentication - logout subcommand
// config - global program config
// args - arguments parsed to subcommand
func authLogoutCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli auth logout
kcmscli auth logout - ends the authenticated session with k8s-cms API.
OPTIONS
`
	// parse & evaluate options
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// logout - remove refreshToken
	config := readConfigFile()
	config.RefreshToken = ""
	config.commit()
}

// authentication - auth check subcommand
// config - global program config
// args - arguments parsed to subcommand
func authCheckCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli auth check
kcmscli auth check - check if the CLI is currently authenticated with the API
OPTIONS
`
	// parse & evaluate options
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// check validity of refresh token with server
	configFile := readConfigFile()
	api := makeAPI(configFile, globalConfig)
	resp := api.call("GET", "auth/check", "", nil)

	// parse server response
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Authenticated with the API")
	case http.StatusUnauthorized:
		die(fmt.Sprint("Not authenticated with the API"))
	default:
		die(fmt.Sprintf("Got unknown status code: %d", resp.StatusCode))
	}
}
