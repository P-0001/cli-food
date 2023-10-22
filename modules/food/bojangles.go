package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartBojanglesTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createBojanglesAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account Created")
	} else {
		task.Logger.Error("Account Creation Failed")
	}

	return task
}

func createBojanglesAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{
		"Host":            []string{"offers-prd--bojangles-dev.netlify.app"},
		"Content-Type":    []string{"application/json"},
		"accept":          []string{"application/json, text/plain, */*"},
		"client_type":     []string{"ios"},
		"version":         []string{"2"},
		"path":            []string{"members/create"},
		"accept-language": []string{"en-us"},
		"User-Agent":      []string{"Bojangles/1 CFNetwork/978.0.7 Darwin/18.7.0"},
		"client_id":       []string{"XMzXxrjALeejfUD2Komc"}, //cSpell: disable-line
	}

	requestData := map[string]any{
		"email":                            task.Account.Email,
		"first_name":                       task.Account.FirstName,
		"last_name":                        task.Account.LastName,
		"password":                         task.Account.Password,
		"phone":                            task.Account.Phone,
		"zip":                              utils.RndZip(),
		"opt_in":                           false,
		"optin_sms":                        false,
		"firebase_push_notification_token": "",
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://offers-prd--bojangles-dev.netlify.app/api",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "token")
}
