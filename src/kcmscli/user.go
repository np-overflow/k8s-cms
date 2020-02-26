/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * user subcommand
 */
package main

import (
	"fmt"
	"github.com/np-overflow/k8s-cms/src/kcmscli/utils"
	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// user type
type User struct {
	UserName string `yaml:"username"`
	FullName string `yaml:"last_name"`
	Passwd   string `yaml:"password"`
}

func getUserHeaders(seperator string) []byte {
	headers := []string{"username", "fullname", "password"}
	return []byte(strings.Join(headers, ","))
}

func (user User) toCSV(seperator string) []byte {
	contents := []string{user.UserName, user.FullName, user.Passwd}
	return []byte(strings.Join(contents, ","))
}

// write users to disk as CSV
func writeUsersCSV(path string, users []User) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	_, err = file.Write(getUserHeaders(","))
	_, err = file.WriteString("\n")
	if err != nil {
		die(err.Error())
	}

	for _, user := range users {
		_, err = file.Write(user.toCSV(","))
		_, err = file.WriteString("\n")
		if err != nil {
			die(err.Error())
		}
	}
}

// write users to disk as YAML
func writeUsers(users []User, contestDir string) {
	usersYAML, err := yaml.Marshal(&users)
	if err != nil {
		die(err.Error())
	}

	// write user yaml to contest dir
	path := filepath.Join(contestDir, "users.yaml")
	utils.WriteBytes(usersYAML, path)
}

// read users form disk as YAML
func readUsers(contestDir string) []User {
	// read user yaml
	path := filepath.Join(contestDir, "users.yaml")
	usersYAML := utils.ReadBytes(path)

	// parse user yaml
	var users []User
	yaml.Unmarshal(usersYAML, &users)

	return users
}

// user subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func userCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli user [options] <subcommand ...>
import, export users

SUBCOMMANDS
import - import users from registration CSV to contest data
export - export contest users to CSV

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
	args = optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <subcommand>")
	}
	subCmd, _ := args[0], args[1:]

	// delegate to various subcommands
	switch subCmd {
	case "import":
		userImportCmd(globalConfig, args)
	case "export":
		userExportCmd(globalConfig, args)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

/* user import */
// user import config
type UserImportConfig struct {
	isPasswdAutogen  bool
	autogenPasswdLen int
	passwdCol        string
	isVerbose        bool
	userNameCol      string
	fullNameCol      string
	registryPath     string
	contestDir       string
	csvSeperator     string
}

// user import subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func userImportCmd(globalConfig *GlobalConfig, args []string) {
	usageInfo := `Usage: kcmscli user import [options] <registration csv> <contest dir>
Registers users by generating user credentials into users.yaml in the given <contest dir> 
with the users found in given registration CSV file.

NOTES:
- The registration CSV's first line should contain column header names

OPTIONS
	`
	// parse command line options
	userNameCol, passwdCol := "user", "AUTO"
	csvSeperator, autogenPasswdLen, fullNameCol := ",", 16, "name"
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.FlagLong(&globalConfig.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.FlagLong(&userNameCol, "user-column", 'u', "Name of the column to "+
		"extract usernames in registration CSV")
	optSet.FlagLong(&passwdCol, "password", 'p', "The password assigned to the users or "+
		"'AUTO' to autogenerate password")
	optSet.FlagLong(&csvSeperator, "csv-seperator", 's', "The seperator used to seperate the CSV file")
	optSet.FlagLong(&autogenPasswdLen, "autogen-len", 'l', "Length of the "+
		"autogenerated password if autogeneration is enabled")
	optSet.FlagLong(&fullNameCol, "name-column", 'n', "Name of the column to "+
		"extract participant's fullname in registration CSV")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	isPasswdAutogen := false
	if passwdCol == "AUTO" {
		isPasswdAutogen = true
		passwdCol = ""
	}

	// parse positional arguments
	args = optSet.Args()
	if len(args) != 2 {
		die("Missing positional arguments: <registration csv> <contest dir>")
	}
	registryPath, contestDir := args[0], args[1]

	if globalConfig.isVerbose {
		fmt.Printf("using username column: \"%s\"\n", userNameCol)
		fmt.Printf("using registration CSV path: \"%s\"\n", registryPath)
		fmt.Printf("using contest directory path: \"%s\"\n", contestDir)
		if isPasswdAutogen {
			fmt.Println("configured to autogenerate passwords for users")
		} else {
			fmt.Printf("using password column: \"%s\"\n", passwdCol)
		}
	}

	// peform user import
	config := UserImportConfig{
		isVerbose:        globalConfig.isVerbose,
		isPasswdAutogen:  isPasswdAutogen,
		passwdCol:        passwdCol,
		userNameCol:      userNameCol,
		fullNameCol:      fullNameCol,
		registryPath:     registryPath,
		contestDir:       contestDir,
		csvSeperator:     csvSeperator,
		autogenPasswdLen: autogenPasswdLen,
	}

	config.importUsers()
}

// validate user import config state
func (config UserImportConfig) validate() {
	exists := func(path string) bool {
		if _, err := os.Stat(path); err == nil {
			return true
		} else {
			return false
		}
	}

	if !config.isPasswdAutogen && len(config.passwdCol) < 1 {
		die("Password column cannot be empty if not using password autogeneration.")
	} else if len(config.userNameCol) < 1 {
		die("Username column cannot be empty")
	} else if !exists(config.registryPath) {
		die("Could not open registration CSV path")
	} else if !exists(config.contestDir) {
		die("Could not open contest directory")
	} else if len(config.csvSeperator) < 1 {
		die("CSV seperator cannot be empty")
	} else if config.autogenPasswdLen < 1 {
		die("Length of autogenerated password cannot be 0 or negative")
	}
}

// compile passwds for the given nUsers.
func (config UserImportConfig) compilePasswds(nUsers int) []string {
	// autogenerate a alpha numeric passwd
	genPasswd := func(genLen int) string {
		alphaNumChars := "abcdefghijklmnopqrstunwxyzABCDEFGHIJKLMNOPQRSTUNWXYZ0123456789"
		var passwdChars []string

		for i := 0; i < genLen; i++ {
			randIdx := int(rand.Uint32()) % len(alphaNumChars)
			passwdChars = append(passwdChars, string(alphaNumChars[randIdx]))
		}

		return strings.Join(passwdChars, "")
	}

	// compile autogenerate passwds
	var passwds []string
	for i := 0; i < nUsers; i++ {
		passwds = append(passwds, genPasswd(config.autogenPasswdLen))
	}

	return passwds
}

// build & returns users for importConfig
func (config UserImportConfig) buildUsers() []User {
	// compile user data
	userNames := utils.LoadCSVColumn(config.registryPath, config.userNameCol, config.csvSeperator)
	fullNames := utils.LoadCSVColumn(config.registryPath, config.fullNameCol, config.csvSeperator)
	nUsers := len(userNames)
	var passwds []string
	if config.isPasswdAutogen {
		passwds = config.compilePasswds(len(userNames))
	} else {
		passwds = utils.LoadCSVColumn(config.registryPath, config.passwdCol, config.csvSeperator)
	}

	// build user structs
	var users []User
	for i := 0; i < nUsers; i++ {
		user := User{
			UserName: userNames[i],
			FullName: fullNames[i],
			Passwd:   passwds[i],
		}
		users = append(users, user)
	}

	return users
}

// perform user imoport based on user import config
func (config UserImportConfig) importUsers() {
	// compile imported/autogenerate user data
	config.validate()
	users := config.buildUsers()
	writeUsers(users, config.contestDir)
}

/* user export */
// user export subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func userExportCmd(globalConfig *GlobalConfig, args []string) {
	usageInfo := `Usage: kcmscli user export  [options] <contest dir>
export contest users to CSV

OPTIONS
`
	// parse & evaluate options
	exportPath := "users.csv"
	optSet := getopt.New()
	optSet.FlagLong(&globalConfig.shouldHelp, "help", 'h', "show usage info")
	optSet.FlagLong(&globalConfig.isVerbose, "verbose", 'v', "produce verbose output")
	optSet.FlagLong(&exportPath, "export-path", 'p', "path of CSV to export user data")
	optSet.Parse(args)

	if globalConfig.shouldHelp {
		fmt.Print(usageInfo)
		optSet.PrintOptions(os.Stdout)
		os.Exit(0)
	}

	// parse positional arguments
	args = optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <contest dir>")
	}
	contestDir := args[0]

	// read users & write as csv
	users := readUsers(contestDir)
	writeUsersCSV(exportPath, users)
}