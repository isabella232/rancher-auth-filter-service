package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
)

//RequestData is for the JSON output
type RequestData struct {
	Headers map[string][]string    `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

//AuthorizeData is for the JSON output
type AuthorizeData struct {
	Message string `json:"message,omitempty"`
}

//MessageData is for the JSON output
type MessageData struct {
	Data []interface{} `json:"data,omitempty"`
}

//ValidationHandler is a handler for cookie token and returns the request headers and accountid and projectid
func ValidationHandler(w http.ResponseWriter, r *http.Request) {

	reqestData := RequestData{}
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Info(err)
		logrus.Infof("Cannot extract the request")
		w.WriteHeader(400)
		return
	}
	praseCookieErr := json.Unmarshal(input, &reqestData)
	if praseCookieErr != nil {
		logrus.Info(praseCookieErr)
		logrus.Infof("Cannot parse the request.")
		w.WriteHeader(400)
		return
	}
	var cookie []string
	if len(reqestData.Headers["Cookie"]) >= 1 {
		cookie = reqestData.Headers["Cookie"]
	} else {
		logrus.Infof("No Cookie found.")
		w.WriteHeader(http.StatusOK)
		return
	}
	var cookieString string
	if len(cookie) >= 1 {

		for i := range cookie {
			if strings.Contains(cookie[i], "token") {
				cookieString = cookie[i]
			}
		}

	} else {
		logrus.Infof("No token found in cookie.")
		w.WriteHeader(http.StatusOK)
		return
	}

	tokens := strings.Split(cookieString, ";")
	tokenValue := ""
	if len(tokens) >= 1 {
		for i := range tokens {
			if strings.Contains(tokens[i], "token") {
				if len(strings.Split(tokens[i], "=")) > 1 {
					tokenValue = strings.Split(tokens[i], "=")[1]
				}

			}

		}
	} else {
		logrus.Infof("Cannot split cookie string")
		w.WriteHeader(http.StatusOK)
		return
	}
	if tokenValue == "" {
		logrus.Infof("No token found")
		w.WriteHeader(http.StatusOK)
		return
	}

	//check if the token value is empty or not
	if tokenValue != "" {
		logrus.Infof("token:" + tokenValue)
		accountID := getValue(manager.URL, "accounts", tokenValue)
		projectID := getValue(manager.URL, "projects", tokenValue)
		//check if the accountID or projectID is empty
		if accountID[0] != "" && projectID[0] != "" {
			if accountID[0] == "Unauthorized" || projectID[0] == "Unauthorized" {
				w.WriteHeader(401)
				logrus.Infof("Token is not valid." + tokenValue)
			} else if accountID[0] == "ID_NOT_FIND" || projectID[0] == "ID_NOT_FIND" {
				w.WriteHeader(501)
				logrus.Infof("Cannot provide the service. Please check the rancher server URL." + manager.URL)
			} else {
				//construct the responseBody
				var headerBody = make(map[string][]string)
				var Body = make(map[string]interface{})

				requestHeader := reqestData.Headers
				for k, v := range requestHeader {
					headerBody[k] = v
				}
				requestBody := reqestData.Body
				for k, v := range requestBody {
					Body[k] = v
				}
				headerBody["X-API-Project-Id"] = projectID
				headerBody["X-API-Account-Id"] = accountID
				var responseBody RequestData
				responseBody.Headers = headerBody
				responseBody.Body = Body
				//convert the map to JSON format
				if responseBodyString, err := json.Marshal(responseBody); err != nil {
					logrus.Info(err)
					w.WriteHeader(500)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(responseBodyString)
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
		logrus.Infof("Cannot connect to the rancher server. Please check the rancher server URL")
		result = []string{"ID_NOT_FIND"}
		return result
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	authMessage := AuthorizeData{}
	err = json.Unmarshal(bodyText, &authMessage)
	if err != nil {
		logrus.Info(err)
		logrus.Infof("Cannot parse the authorization data.")
		result = []string{"ID_NOT_FIND"}
		return result
	}
	if authMessage.Message == "Unauthorized" {
		result = []string{"Unauthorized"}
	} else {
		messageData := MessageData{}
		err = json.Unmarshal(bodyText, &messageData)
		if err != nil {
			logrus.Info(err)
			logrus.Infof("Cannot parse the id.")
			result = []string{"ID_NOT_FIND"}

		}
		//get id from the data

		for i := 0; i < len(messageData.Data); i++ {

			idData, suc := messageData.Data[i].(map[string]interface{})
			if suc {
				id, suc := idData["id"].(string)
				name, namesuc := idData["name"].(string)
				if suc && namesuc {
					result = append(result, id)
					//if the token belongs to admin, only return the admin token
					if name == "admin" && path == "accounts" {
						result = []string{id}
						break
					}
				} else {
					logrus.Infof("No id find")
					result = []string{"ID_NOT_FIND"}
				}
			}

		}
		//get the admin user id. admin token will list all the ids. Need to just keep admin id.

	}

	return result
}
