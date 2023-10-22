package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

var tacoBellVersion string = "8.26.0"

func (task FoodTaskS) StartTacoBellTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createTacoBellAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account Created")
	} else {
		task.Logger.Error("Account Creation Failed")
	}

	return task
}

func createTacoBellAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}
	headers.Set("Host", "www.tacobell.com")
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Datadog-Parent-ID", "")
	headers.Set("X-Datadog-Sampled", "1")
	headers.Set("Accept", "*/*")
	headers.Set("X-Datadog-Sampling-Priority", "1")
	headers.Set("X-Datadog-Trace-ID", "")
	headers.Set("X-ACF-Sensor-Data", "")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("User-Agent", "Taco Bell/"+tacoBellVersion+"-iOS")
	headers.Set("X-Datadog-Origin", "rum")

	requestData := map[string]any{
		"defaultNotification": false,
		"lastName":            task.Account.LastName,
		"birthday":            utils.Birthday("2000-{M}-{D}", false, true),
		"uid":                 task.Account.Email,
		"subscription":        false,
		"firstName":           task.Account.FirstName,
		"password":            task.Account.Password,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://www.tacobell.com/tacobellwebservices/v2/tacobell/users",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 201,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	strRes := string(res)

	return strings.Contains(strRes, "customerID")
}
