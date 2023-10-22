package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

const wendysVERSION = "10.0.10"

func (task FoodTaskS) StartWendysTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     "",
		Password:  utils.SecurePassword("") + "W1w!",
	}

	res := createWendysAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createWendysAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}
	ua := fmt.Sprintf("Wendys/%s (iPhone; iOS %s; Scale/3.00)", wendysVERSION, utils.RndIosVersion3())

	headers.Set("Host", "customerservices.wendys.com")
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", ua)
	headers.Set("Accept-Language", "en-US;q=1")

	requestData := map[string]any{
		"lastName":            task.Account.LastName,
		"firstName":           task.Account.FirstName,
		"deviceId":            utils.DeviceUUID(true, false, true),
		"optIn":               false,
		"postal":              utils.RndZip(),
		"hasLoyalty":          false,
		"collectSecurityCode": false,
		"password":            task.Account.Password,
		"login":               task.Account.Email,
		"currency":            "USD",
		"accountStatus":       0,
		"hasLoyalty2020":      false,
		"cntry":               "US",
		"birthdate":           utils.Birthday("{M}{D}1904", false, true),
		"hasPasscode":         false,
		"terms":               true,
		"lang":                "en",
		"isLoyaltyProfile":    true,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	url := "https://customerservices.wendys.com/CustomerServices/rest/createProfile?lang=en&cntry=US&sourceCode=MY_WENDYS&version=" + wendysVERSION

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              url,
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "SUCCESS")
}
