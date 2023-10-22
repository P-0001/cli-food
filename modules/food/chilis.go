package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartChillisTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.RandomRealPhoneNumber(),
		Password:  utils.SecurePassword(""),
	}

	res := createChillisAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func getOauth(task FoodTaskS, combyte string) (bool, string) {
	headers := client.TLSHeaders{}

	headers.Set("Host", "gsapi.brinker.com")
	headers.Set("Accept", "*/*")
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("combyte", combyte)
	headers.Set("User-Agent", "Chili's/2.0 CFNetwork/978.0.7 Darwin/18.7.0")
	headers.Set("Accept-Language", "en-us")
	headers.Set("Authorization", "Basic Y2hpbGlzX2lvczphb0FzbjR4S29mQk9KWGxpQjZEanFLTWRHREE0WTAzclQwSFlmQ3c4")

	form := utils.ToForm(map[any]any{
		"grant_type": "client_credentials",
	})

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://gsapi.brinker.com/v1.0.0/oauth/token",
		Headers:          headers,
		Body:             bytes.NewReader([]byte(form)),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false, ""
	}

	var data map[string]any

	err = json.Unmarshal(res, &data)

	if err != nil {
		task.Logger.Error("Json Unmarshal Failed", err.Error())
		return false, ""
	}

	auth := data["access_token"].(string)

	return true, auth
}

func addPhone(task FoodTaskS, combyte string, auth string, which int) bool {
	headers := client.TLSHeaders{}

	headers.Set("Host", "gsapi.brinker.com")
	headers.Set("Accept", "*/*")
	headers.Set("Content-Type", "application/json")
	headers.Set("X-NewRelic-ID", "")
	headers.Set("Authorization", "Bearer "+auth)
	headers.Set("combyte", combyte)
	headers.Set("User-Agent", "Chili's/2.0 CFNetwork/978.0.7 Darwin/18.7.0")
	headers.Set("Accept-Language", "en-us")

	var requestData any

	if which == 1 {
		requestData = map[string]any{
			"phoneNumber": task.Account.Phone,
		}
	} else {
		requestData = map[string]any{
			"phoneNumber": task.Account.Phone,
			"email":       task.Account.Email,
			"userState":   -1,
			"channelType": 1,
		}
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://gsapi.brinker.com/v1.0.0/api/join/verify/phoneEmail",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	_ = string(res)

	return true
}

func createChillisAccount(task FoodTaskS) bool {

	combyte := utils.DeviceUUID(true, false, false)

	succ, auth := getOauth(task, combyte)

	if !succ {
		task.Logger.Error("Failed to get oauth")
		return false
	}

	// println("Auth :" + auth)

	succ = addPhone(task, combyte, auth, 1)

	if !succ {
		task.Logger.Error("Failed to add phone 1")
		return false
	}

	succ = addPhone(task, combyte, auth, 2)

	if !succ {
		task.Logger.Error("Failed to add phone 2")
		return false
	}

	headers := client.TLSHeaders{}

	headers.Add("Host", "gsapi.brinker.com")
	headers.Add("Accept", "*/*")
	headers.Add("Content-Type", "application/json")
	headers.Add("X-NewRelic-ID", "")
	headers.Add("Authorization", "Bearer "+auth)
	headers.Add("combyte", "")
	headers.Add("User-Agent", "Chili's/2.0 CFNetwork/978.0.7 Darwin/18.7.0")
	headers.Add("Accept-Language", "en-us")

	requestData := map[string]any{

		"over18":            true,
		"userState":         0,
		"mobileOptIn":       false,
		"email":             task.Account.Email,
		"enrollmentChannel": nil,
		"dob":               utils.Birthday("{Y}-{M}-{D}", false, true),
		"firstName":         task.Account.FirstName,
		"enrollmentDate":    "",
		"storeGUID":         "",
		"channelType":       2,
		"terminalID":        "",
		"storeCode":         "0010050960",
		"source":            "chilisapp",
		"phoneNumber":       task.Account.Phone,
		"channelID":         "",
		"password":          task.Account.Password,
		"lastName":          task.Account.LastName,
	}

	requestData["address"] = map[string]any{
		"zip":     utils.RndZip(),
		"country": "US",
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://services-mob.panerabread.com/register",
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
