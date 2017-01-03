package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/rancher/rancher-auth-filter-service/manager"
)

//ValidationHandler is a handler for cookie token and returns the request headers and accountid and projectid
func ValidationHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == nil {
		//check if the token value is empty or not
		if cookie.Value != "" {
			logrus.Infof("token:" + cookie.Value)
			accountID := getValue(manager.Url, "accounts", cookie.Value)
			projectID := getValue(manager.Url, "projects", cookie.Value)
			//check if the accountID or projectID is empty
			if accountID[0] != "" && projectID[0] != "" {
				if accountID[0] == "Unauthorized" || projectID[0] == "Unauthorized" {
					w.WriteHeader(401)
					logrus.Infof("Token is not valid." + cookie.Value)
				} else if accountID[0] == "ID_NOT_FIND" || projectID[0] == "ID_NOT_FIND" {
					w.WriteHeader(501)
					logrus.Infof("Cannot provide the service. Please check the rancher server URL." + manager.Url)
				} else {
					//construct the responseBody
					var headerBody map[string][]string = make(map[string][]string)
					for k, v := range r.Header {
						headerBody[k] = v
					}
					headerBody["X-API-Project-Id"] = projectID
					headerBody["X-API-Account-Id"] = accountID
					var responseBody map[string]map[string][]string = make(map[string]map[string][]string)
					responseBody["headers"] = headerBody
					//convert the map to JSON format
					if responseBodyString, err := json.Marshal(responseBody); err != nil {
						panic(err)
					} else {
						fmt.Fprintln(w, string(responseBodyString))
					}
				}
			}

		}

	}
}

//get the projectID and accountID from rancher API
func getValue(host string, path string, token string) []string {
	var result []string
	client := &http.Client{}
	requestURL := host + "v2-beta/" + path
	req, err := http.NewRequest("GET", requestURL, nil)
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	js, _ := simplejson.NewJson(bodyText)
	authorized, _ := js.Get("message").String()

	if authorized == "Unauthorized" {
		result = []string{"Unauthorized"}
	} else {
		var id string
		jsonBody, _ := simplejson.NewJson(bodyText)
		dataLenth := len(jsonBody.Get("data").MustArray())
		for i := 0; i < dataLenth; i++ {
			id, err = jsonBody.Get("data").GetIndex(i).Get("id").String()

			if err != nil {
				logrus.Info(err)
				result = []string{"ID_NOT_FIND"}
			} else {
				result = append(result, id)
			}
		}

	}

	return result
}
