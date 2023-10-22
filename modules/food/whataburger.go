package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartWhataburgerTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.RandomRealPhoneNumber(),
		Password:  utils.SecurePassword(""),
	}

	res := createWhataburgerAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createWhataburgerAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}
	headers.Set("content-type", "application/json")
	headers.Set("authority", "api.whataburger.com")
	headers.Set("accept", "application/json")
	headers.Set("accept-language", "en-US,en;q=0.9")
	headers.Set("origin", "https://whataburger.com")
	headers.Set("referer", "https://whataburger.com/")
	headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	headers.Set("x-api-key", "E08F3550-23FE-4360-BD6C-08314E6C3E2F")
	headers.Set("x-client", "SPA")
	headers.Set("x-device-fingerprint", utils.DeviceUUID(false, false, false))
	headers.Set("x-device-id", utils.DeviceUUID(true, false, false))

	requestData := map[string]any{
		"firstname":           task.Account.FirstName,
		"lastname":            task.Account.LastName,
		"email":               task.Account.Email,
		"zipcode":             utils.RndZip(),
		"password":            task.Account.Password,
		"hasSubcribedToEmail": false,
		"phoneNumber":         task.Account.Phone,
		"confPW":              task.Account.Password,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://api.whataburger.com/v2.4/accounts/signup",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "accessToken")
}
