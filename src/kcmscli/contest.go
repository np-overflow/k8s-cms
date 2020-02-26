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

// push contest with the specified contest data to api
// return the response from the api push request
func pushContest(globalConfig *GlobalConfig, method string, contestDir string) *http.Response {
	// package contest into an archive to send
	if globalConfig.isVerbose {
		fmt.Println("Packing contest for upload...")
	}
	workDir := utils.MakeTempDir("archive")
	archivePath := filepath.Join(workDir, "contest.tgz")
	packageContest(contestDir, archivePath)
	defer os.RemoveAll(workDir)

	// construct form data with archive to send with request
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
	if globalConfig.isVerbose {
		fmt.Println("Pushing contest to API...")
	}
	api := makeAPI(readConfigFile(), globalConfig)
	api.refreshAccess() 
	resp := api.call(method, "contest/import", contentType, formdata)

	return resp
}

// contest  subcommand
// globalConfig - global program config
// args - arguments parsed to subcommand
func contestCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest [options] <subcommand ...>
Import, List, Delete contests

SUBCOMMANDS
import - import a contest into k8s-cms
update - update a contest in k8s-cms
list - list contests in k8s-cms 
get - get contest infomation
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
	case "get":
		contestGetCmd(globalConfig, args)
	case "update":
		contestUpdateCmd(globalConfig, args)
	case "delete":
		contestDeleteCmd(globalConfig, args)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

/* contest CRUD operations */
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
	
	// make api request
	resp := pushContest(globalConfig, "POST", contestDir)
	
	// parse results of API call
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Imported contest")
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	case http.StatusConflict:
		die("Attempted to import a duplicate contest")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", resp.StatusCode))
	}
}

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

// get contests informations
// config - global program config
// args - arguments parsed to subcommand
func contestGetCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest get <contest id>
Get contest info with the contest with specific id.

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
	
	// parse positional arguments
	args = optSet.Args()
	if len(args) < 1 {
		die("Missing positional arguments: <contest id>")
	}
	contestId := args[0]

	// perform api call & read response
	api := makeAPI(readConfigFile(), globalConfig)
	api.refreshAccess()
	statusCode, contestInfo := api.callJSON("GET", "contest/" + contestId, nil)
	
	// parse response code
	switch statusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		die("Contest with given id could not be found")
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", statusCode))
	}

	// render contest info
	fmt.Printf("CONTEST %s\n", contestInfo["name"])
	fmt.Printf("description: %s\n", contestInfo["description"])
	fmt.Printf("localizations: %s\n", contestInfo["allowedLocalizations"])
	fmt.Printf("languages: %s\n", contestInfo["languages"])
	fmt.Printf("timezone: %s\n", contestInfo["timezone"])
	fmt.Printf("start time: %s\n", contestInfo["startTime"])
	fmt.Printf("end time: %s\n", contestInfo["stopTime"])
	fmt.Printf("max contest time: %.2fs\n", contestInfo["maxUserContestTime"])
	fmt.Printf("max submissions: %d\n", int(contestInfo["maxSubmissionNum"].(float64)))

	var enabled []string
	var disabled []string
	features := []string{ 
		"allowIPAutoLogin", "allowPasswordAuthentication", "allowQuestions",
		"allowSubmissionsDownload", "enforceIPRestriction",
	}
	
	for _, feature := range features {
		if contestInfo[feature].(bool) {
			enabled = append(enabled, feature)
		} else {
			disabled = append(disabled, feature)
		}
	}
	
	fmt.Printf("enabled: %s\n", enabled)
	fmt.Printf("disabled: %s\n", disabled)
}

// contest update subcommand subcommand
func contestUpdateCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest update [options] <contest dir>
Update contest into k8s-cms using contest data at contest dir. 
The contest data should be in the italy contest format https://github.com/cms-dev/con_test.

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
	
	// make api request
	resp := pushContest(globalConfig, "PATCH", contestDir)
	
	// parse results of API call
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Updated contest")
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", resp.StatusCode))
	}
}

func contestDeleteCmd(globalConfig *GlobalConfig, args []string) {
	var usageInfo string = `Usage: kcmscli contest delete [options] <contest id>
Delete contest info with the contest with specific id.

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
		die("Missing positional arguments: <contest id>")
	}
	contestId := args[0]

	// make api call to delete contest
	api := makeAPI(readConfigFile(), globalConfig)
	api.refreshAccess()
	statusCode, _ := api.callJSON("DELETE", "contest/" + contestId, nil)
	
	switch statusCode {
	case http.StatusOK:
		fmt.Println("Deleted contest")
	case http.StatusUnauthorized:
		die("Not authorized to perform command: Login first. ")
	default:
		die(fmt.Sprintf("Got unknown status code: %d", statusCode))
	}
}
