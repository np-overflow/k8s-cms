/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * user subcommand
 */
package main

import (
	"fmt"
	"net/http"
	"github.com/np-overflow/k8s-cms/src/kcmscli/utils"
	"github.com/otiai10/copy"
	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

// contest type
type Contest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	TokenMode   string `yaml:"token_mode"`
	Timezone    string `yaml:"timezone"`

	Start int `yaml:"start"`
	End   int `yaml:"end"`

	PerUserTime         int `yaml:"per_user_time"`
	MaxSubmitCount      int `yaml:"max_submission_number"`
	MaxUserTestCount    int `yaml:"max_user_test_number;"`
	MinSubmitInterval   int `yaml:"min_submission_interval"`
	MinUserTestInterval int `yaml:"min_user_test_interval"`

	ProgrammingLanguages []string `yaml:"programming_languages"`
	Tasks                []string `yaml:"tasks"`

	Users []User `yaml:"users"`
}

// init contest with sane defaults
func NewContest() Contest {
	contest := Contest{
		Start:               0,
		End:                 1,
		PerUserTime:         18000,
		MinSubmitInterval:   20,
		MinUserTestInterval: 20,
		MaxSubmitCount:      99999,
		MaxUserTestCount:    99999,
		ProgrammingLanguages: []string{
			"C++11 / g++",
		},
	}

	return contest
}

// read contest metadata data from the given contest dir into contest struct
// reads & merges user data from user.yaml into contest struct
// returns the contest with metadata
func readMetadata(contestDir string) Contest {
	contest := NewContest()

	// merge users in users.yaml into contest users in contest.yml
	metaPath := filepath.Join(contestDir, "contest.yaml")
	err := yaml.Unmarshal(utils.ReadBytes(metaPath), &contest)
	if err != nil {
		die(err.Error())
	}
	contest.Users = append(contest.Users, readUsers(contestDir)...)

	return contest
}

// package the given contest into an single archive file suitable
// for network transmission. writes the archive to the given archive path
func packageContest(contestDir string, archivePath string) {
	// copy contest dir into as working contest data
	workDir := utils.MakeTempDir("package_workspace")
	defer os.RemoveAll(workDir)
	workContestDir := filepath.Join(workDir, filepath.Base(contestDir))
	copy.Copy(contestDir, workContestDir)

	// read contest metadata and and rewrite metadata into work directory
	contest := readMetadata(contestDir)
	contestYAML, err := yaml.Marshal(contest)
	if err != nil {
		die(err.Error())
	}
	metaPath := filepath.Join(workContestDir, "contest.yaml")
	utils.WriteBytes(contestYAML, metaPath)

	// create tar archive of contest data
	utils.MakeTGZ(workContestDir, archivePath)
}

// contest  subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func contestCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest [options] <subcommand ...>
Import, List, Delete contests

SUBCOMMANDS
import - import a contest into k8s-cms
list - list contests in k8s-cms 
delete - delete a contest 

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
	subCmd := args[0]

	// delegate to various subcommands
	switch subCmd {
	case "import":
		contestImportCmd(globalConfig, args)
	case "list":
		contestListCmd(globalConfig, args)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

/* import contests */
// contest import subcommand
func contestImportCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest import [options] <contest dir>
Import contest into k8s-cms using contest data at contest dir. 
The contest data should be in the italy contest format https://github.com/cms-dev/con_test.
Users should already be imported into contest dir with kcmscli user import

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

	// parse position args
	args = optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <contest dir>")
	}
	contestDir := args[0]
	if _, err := os.Stat(contestDir); os.IsNotExist(err) {
		die("Given contest directory does not exist")
	}

	// package contest into an archive to send
	workDir := utils.MakeTempDir("archive")
	archivePath := filepath.Join(workDir, "contest.tgz")
	packageContest(contestDir, archivePath)
	defer os.RemoveAll(workDir)

	// construct form data with archive to send with request
	fmt.Println(archivePath)
	archiveFile, err := os.Open(archivePath)
    if err != nil {
        panic(err.Error())
    }
	defer archiveFile.Close()
	fileFields := map[string]*os.File {
		"file": archiveFile,
	}
	contentType, formdata := utils.NewMultipartData(map[string]string{}, fileFields)
	
	// make api call to import contes 
	api := makeAPI(readConfigFile(), globalConfig)
	api.refreshAccess() 
	resp := api.call("POST", "contest/import", contentType, formdata)

	// parse results of API call
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("Imported contest")
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	case http.StatusConflict:
		die("Attempted to import a duplicate contest")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", resp.StatusCode))
	}
}

/* list & delete contests */
// list contests subcommand
// config - global program config
// args - arguments parsed to subcommand
func contestListCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest list [options]
List available contests in k8s-cms

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
	
	// perform api call & read response
	api := makeAPI(readConfigFile(), globalConfig)
	api.refreshAccess()
	// include names in the response 
	params := url.Values{}
	params.Set("incl-names", "1")
	resp := api.call("GET", "contests?" + params.Encode(), "", nil)
	responseJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		die(err.Error())
	}

	// parse results of API call
	var contestInfo []map[string]interface{}
	switch resp.StatusCode {
	case http.StatusOK:
		json.Unmarshal(responseJSON, &contestInfo)
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", resp.StatusCode))
	}
	
	// render contest info
	fmt.Println("ID\tCONTEST")
	for _, info := range contestInfo {
		id := int(info["id"].(float64))
		fmt.Printf("%d\t%s\n", id, info["name"])
	}
}
