package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

var dqVersion string = "3.1.19"

func (task FoodTaskS) StartDQTask() FoodTaskS {

	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createDQAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createDQAccount(task FoodTaskS) bool {
	os := utils.RndIosVersion3()
	headers := client.TLSHeaders{}

	headers.Set("Host", "prod-api.dairyqueen.com")
	headers.Set("accept", "*/*")
	headers.Set("content-type", "application/json")
	headers.Set("country", "US")
	headers.Set("partner-platform", "iOS")
	headers.Set("x-device-id", utils.DeviceUUID(true, false, true))
	//cSpell:disable-next-line
	headers.Set("user-agent", fmt.Sprintf("DairyQueen/%s (com.dairyqueen.ios.loyaltyapp.production; build:691; iOS %s) Alamofire/5.4.3", dqVersion, os))
	headers.Set("accept-language", "en-US;q=1.0")

	payload := fmt.Sprintf(`
		mutation {
			signup(
				userInput: {
					firstName: "%s"
					lastName: "%s"
					emailAddress: "%s"
					phoneNumber: "%s"
					password: "%s"
					passwordConfirmation: "%s"
					marketingPnSubscription: false
					subscribeToMarketingEmails: false
					termsConditions: true
					favoriteLocations: "334482"
					birthday: "%s"
				}
			) {
				id
			}
		}
	`, task.Account.FirstName, task.Account.LastName, task.Account.Email, task.Account.Phone, task.Account.Password, task.Account.Password, utils.Birthday("1900-{M}-{D}", false, true))

	requestData := map[string]any{
		"query": payload,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://prod-api.dairyqueen.com/graphql",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), `signup`)
}
