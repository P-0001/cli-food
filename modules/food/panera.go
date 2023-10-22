package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartPaneraTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.RandomRealPhoneNumber(),
		Password:  utils.SecurePassword(""),
	}

	res := createPaneraAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createPaneraAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}

	headers.Set("Host", "services-mob.panerabread.com")
	headers.Set("Accept", "*/*")
	headers.Set("appVersion", "4.82.1")
	headers.Set("X-acf-sensor-data", "")
	headers.Set("api_token", "bcf0be75-0de6-4af0-be05-13d7470a85f2")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("deviceId", utils.DeviceUUID(true, false, false))
	headers.Set("User-Agent", "Panera/4.82.1 (iPhone; iOS 15.4.1; Scale/2.00)")
	headers.Set("X-Registration-Source", "Full")
	headers.Set("Content-Type", "application/json")

	month, day, _ := utils.BirthdayObject(false, true)

	requestData := map[string]any{
		"birthDate":    "",
		"emailAddress": task.Account.Email,
		"password":     task.Account.Password,
		"lastName":     task.Account.LastName,
		"phoneNumber":  task.Account.Phone,
		"username":     task.Account.Email,
		"opt":          "false",
		"firstName":    task.Account.FirstName,
	}

	requestData["birthDate"] = map[string]string{
		"birthMonth": month,
		"birthDay":   day,
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
