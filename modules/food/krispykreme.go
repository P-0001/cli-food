package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func (task FoodTaskS) StartKKTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createKKAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createKKAccount(task FoodTaskS) bool {

	os := utils.RndIosVersion()

	ua := fmt.Sprintf("aws-sdk-iOS/2.25.0 iOS/%s en_US", os)

	headers := client.TLSHeaders{
		"Accept":       []string{"application/json"},
		"Host":         []string{"api.krispykreme.com"},
		"Content-Type": []string{"application/json"},
		"User-Agent":   []string{ua},
		"X-D-Token":    []string{""},
	}

	requestData := map[string]any{
		"appVersion":  "22.14.0",
		"email":       task.Account.Email,
		"firstName":   task.Account.FirstName,
		"lastName":    task.Account.LastName,
		"phoneNumber": task.Account.Phone,
		"password":    task.Account.Password,
		"source":      "iOS",
		"zipCode":     utils.RndZip(),
		"birthday":    utils.Birthday("{M}/{D}/2000", false, true),
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://api.krispykreme.com/auth/createaccount",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "accessToken")
}
