package food

import (
	"bytes"
	"emerald/client"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartPopeyesTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: strings.ToLower(task.Account.FirstName),
		LastName:  "",
		Email:     strings.ToLower(task.Account.Email),
		Phone:     "",
		Password:  "",
	}

	res := createPopeyesAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createPopeyesAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}

	headers.Add("Host", "use1-prod-plk-gateway.rbictg.com")
	headers.Add("Content-Type", "application/json")
	headers.Add("Accept", "*/*")
	headers.Add("Accept-Language", "en-us")
	headers.Add("User-Agent", "Popeyes/10 CFNetwork/978.0.7 Darwin/18.7.0")
	headers.Add("X-UI-Region", "US")
	headers.Add("X-UI-Language", "en")
	headers.Add("X-Datadog-Origin", "rum")

	requestData := map[string]any{
		"operationName": "SignUp",
		"query":         "mutation SignUp($input: SignUpUserInput!) {\n  signUp(userInfo: $input)\n}",
		"variables": map[string]any{
			"input": map[string]any{
				"country":                "USA",
				"dob":                    "",
				"name":                   task.Account.FirstName,
				"phoneNumber":            "",
				"platform":               "app",
				"stage":                  "prod",
				"userName":               task.Account.Email,
				"wantsPromotionalEmails": false,
				"zipcode":                "",
			},
		},
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://use1-prod-plk-gateway.rbictg.com/graphql",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "signUp")
}
