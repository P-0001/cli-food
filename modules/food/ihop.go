package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

var versionIhop = "4.22.0"

func (task FoodTaskS) StartIhopTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.RandomRealPhoneNumber(),
		Password:  utils.SecurePassword(""),
	}

	res := createIhopAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func getVersion(task FoodTaskS) (bool, string) {
	ua := fmt.Sprintf("ihop/%s/2244(iPhone;iOS;%s;3.0)", versionIhop, utils.RndIosVersion3())

	headers := client.TLSHeaders{}
	headers.Set("Host", "www.ihop.com")
	headers.Set("accept", "application/json")
	headers.Set("content-type", "application/json")
	headers.Set("user-agent", ua)
	headers.Set("accept-language", "en-US,en;q=0.9")

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "GET",
		Url:              "https://www.ihop.com/api/mobile/ios/version",
		Headers:          headers,
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

	v := data["iosversion"].(string)

	return true, v
}

func genToken(task FoodTaskS) (bool, string) {
	headers := client.TLSHeaders{}
	headers.Set("Host", "www.dinefranchisees.com")
	headers.Set("Accept", "*/*")
	headers.Set("X-dine-channel", "IOS")
	headers.Set("User-Agent", "IHOP/2244 CFNetwork/1408.0.4 Darwin/22.5.0")
	headers.Set("Authorization", "Basic MG9hMXN5ZnZqOGdKcXBBeDIwaDg6YUY5bjd4WnNfcF95bVlvblQ1b25CQ243eUhhTkxPUmxUNkxNZEF4RzhETkxFS25MN2J4M0ptQXdFOFpDR3JtSg==")
	headers.Set("Accept-Language", "en-US;q=1.0")
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://www.dinefranchisees.com/oauth2/aus1l1nodzbys0uhX0h8/v1/token?grant_type=client_credentials&scope=read",
		Headers:          headers,
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

func createIhopAccount(task FoodTaskS) bool {

	succ, version := getVersion(task)

	if !succ {
		task.Logger.Error("Failed", "No Version")
		return false
	}

	succ, auth := genToken(task)

	if !succ {
		task.Logger.Error("Failed", "No Auth Token")
		return false
	}

	os := utils.RndIosVersion3()

	ua := fmt.Sprintf("IHOP/%s (com.dineequity.ihop; build:2244; iOS %s) Alamofire/5.0.5", version, os)

	headers := client.TLSHeaders{}
	headers.Set("Host", "mule-api-lb-prod.dinebrands.com")
	headers.Set("x-dine-channel", "IOS")
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "en-US;q=1.0")
	headers.Set("User-Agent", ua)
	headers.Set("Authorization", "Bearer "+auth)

	createdOn := utils.FormatDate("{y}-{m}-{d} {h}:{i}:{s}", utils.GetEstDate(), true)

	type Address struct {
		PostalCode string `json:"postalCode"`
	}

	requestData := map[string]any{
		"marketingEmailSubscription": false,
		"customerPreferredLocation":  "",
		"customerPreferredLocale":    "en/US",
		"lastName":                   task.Account.LastName,
		"emailAddress":               task.Account.Email,
		"firstName":                  task.Account.FirstName,
		"password":                   task.Account.Password,
		"termsandconditions":         true,
		"customerBirthday":           utils.Birthday("{Y}-{M}-{D}", false, true),
		"addresses":                  []Address{{PostalCode: utils.RndZip()}},
		"contactNumber":              utils.ToInt(task.Account.Phone),
		"createdOn":                  createdOn,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://mule-api-lb-prod.dinebrands.com/api/accountservices/v1.0/accounts",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 201,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "accessToken")
}
