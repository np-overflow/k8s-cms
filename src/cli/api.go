/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * API support
*/
package main

import  (
	"io"
	"fmt"
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

/* API */
type API struct {
	globalConfig *GlobalConfig
	refreshToken string
	accessToken string
	apiHost string
	client *http.Client
}

// construct a api object  using the given globalConfig and configFile
func makeAPI(configFile ConfigFile, globalConfig *GlobalConfig) API {
	return API {
		globalConfig: globalConfig,
		refreshToken: configFile.RefreshToken,
		apiHost: configFile.ApiHost,
		client: &http.Client{},
	}
}

// refresh access token, which is required to access specific parts of the API
func (api *API) refreshAccess() error {
	// make api call to refresh access token
	statusCode, response := api.callJSON("GET", "auth/refresh", nil)
	if statusCode != http.StatusOK {
		return errors.New("Failed to refresh access token")
	}
	
	api.accessToken = response["accessToken"]
	return nil
}

// attempt attach a token to authenticate with the API to the given request
func (api API) attachToken(req *http.Request) {
	// attach auth token for authentication
	attachToken := func (token string) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", token))
	}
	if len(api.accessToken) > 0 {
		attachToken(api.accessToken)
		if api.globalConfig.isVerbose {
			fmt.Println("Using access token to authenticate")
		}
	} else if len(api.refreshToken) > 0 {
		attachToken(api.refreshToken)
		if api.globalConfig.isVerbose {
			fmt.Println("Using refresh token to authenticate")
		}
	} else {
		if api.globalConfig.isVerbose {
			fmt.Println("API call will be unauthenticated")
		}
	}
}

// make an API call at the given route with an body, using http method
// calls api at <ApiHost>/api/v0/<route>, where ApiHost is found in the config file.
// returns the response from the http request
func (api API) call(method string, route string, contentType string, body io.Reader) *http.Response {
	// build api call URL
	configFile := readConfigFile()
	apiRoute := fmt.Sprintf("%s/api/v0/%s", configFile.ApiHost, route)

	// build request
	req, err := http.NewRequest(method, apiRoute, body)
	if err != nil {
		die(err.Error())
	}
	req.Header.Add("Content-Type", contentType)
	api.attachToken(req)

	// make call
	if api.globalConfig.isVerbose {
		fmt.Printf("Making API call: %s\n", apiRoute)
	}

	resp, err :=  api.client.Do(req)
	if err != nil {
		die(err.Error())
	}
	
	if api.globalConfig.isVerbose  {
		fmt.Printf("API returned status code: %d\n", resp.StatusCode)
	}
	return resp
}

// make a fully (request, response) JSON API call at the given route, using http method
// returns status code and JSON response as map[string]string
func (api API) callJSON(method string, route string, body interface{}) (
	int, map[string]string) {
	bodyJSON , err := json.Marshal(body)
	if err != nil {
		die(err.Error())
	}
	
	resp := api.call(method, route, "application/json", bytes.NewReader(bodyJSON))
	
	// parse JSON body to map
	var responseJSON []byte
	if body != nil { 
		responseJSON, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			die(err.Error())
		}
	}
	var responseMap map[string]string
	json.Unmarshal(responseJSON, &responseMap)
	
	return resp.StatusCode, responseMap
}
