//
// Written By : @ice3man (Nizamul Rana)
//
// Distributed Under MIT License
// Copyrights (C) 2018 Ice3man
//

// NOTE : We are using Virustotal API here Since we wanted to eliminate the
// rate limiting performed by Virustotal on scraping.
// Direct queries and parsing can be also done :-)

// A Virustotal Client for Subdomain Enumeration
package virustotal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Ice3man543/subfinder/libsubfinder/helper"
)

type virustotalapi_object struct {
	Subdomains []string `json:"subdomains"`
}

var virustotalapi_data virustotalapi_object

// Local function to query virustotal API
// Requires an API key
func queryVirustotalApi(domain string, state *helper.State) (subdomains []string, err error) {

	// Make a search for a domain name and get HTTP Response
	resp, err := helper.GetHTTPResponse("https://www.virustotal.com/vtapi/v2/domain/report?apikey="+state.ConfigState.VirustotalAPIKey+"&domain="+domain, state.Timeout)
	if err != nil {
		return subdomains, err
	}

	// Get the response body
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return subdomains, err
	}

	// Decode the json format
	err = json.Unmarshal([]byte(resp_body), &virustotalapi_data)
	if err != nil {
		return subdomains, err
	}

	// Append each subdomain found to subdomains array
	for _, subdomain := range virustotalapi_data.Subdomains {

		// Fix Wildcard subdomains containg asterisk before them
		if strings.Contains(subdomain, "*.") {
			subdomain = strings.Split(subdomain, "*.")[1]
		}

		if state.Verbose == true {
			if state.Color == true {
				fmt.Printf("\n[%sVIRUSTOTAL%s] %s", helper.Red, helper.Reset, subdomain)
			} else {
				fmt.Printf("\n[VIRUSTOTAL] %s", subdomain)
			}
		}

		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

/*func queryVirustotal(state *helper.State) (subdomains []string, err error) {

	subdomainRegex, err := regexp.Compile("<a target=\"_blank\" href=\"/en/domain/.*\">
      (.*)
    </a>")
	if err != nil {
		return subdomains, err
	}
}*/

// Query function returns all subdomains found using the service.
func Query(args ...interface{}) interface{} {

	domain := args[0].(string)
	state := args[1].(*helper.State)

	var subdomains []string

	// We have recieved an API Key
	// Now, we will use Virustotal API key to fetch subdomain info
	if state.ConfigState.VirustotalAPIKey != "" {

		// Get subdomains via API
		subdomains, err := queryVirustotalApi(domain, state)

		if err != nil {
			fmt.Printf("\nvirustotal: %v\n", err)
			return subdomains
		}
	}

	return subdomains
}
