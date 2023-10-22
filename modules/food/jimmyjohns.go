package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func (task FoodTaskS) StartJimmyjohnsTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createJimmyjohnsAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

type makeResponseJJ struct {
	AuthToken string `json:"authtoken"`
}

func createJimmyjohnsAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}

	headers.Set("Host", "ordering.api.olo.com")
	headers.Set("accept", "application/json")
	headers.Set("content-type", "application/json")
	headers.Set("accept-language", "en-US,en;q=0.9")
	headers.Set("x-device-id", utils.DeviceUUID(true, false, true))
	headers.Set("user-agent", "Jimmy Johns/5.2.0 (bundle id: com.jimmyjohns.jimmyjohns; build: 20269; iOS: 15.5.0)")
	//cSpell: disable-next-line
	headers.Set("authorization", "OloKey RQGfKfsvG5kSw0gDN1jCoC967R7GqECd")

	requestData := map[string]string{
		"contactnumber": task.Account.Phone,
		"firstname":     task.Account.FirstName,
		"lastname":      task.Account.LastName,
		"emailaddress":  task.Account.Email,
		"password":      task.Account.Password,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://ordering.api.olo.com/v1.1/users/create",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	var parsedResponse makeResponseJJ

	// Unmarshal the JSON response into the parsedResponse variable
	err = json.Unmarshal(res, &parsedResponse)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	if parsedResponse.AuthToken == "" {
		task.Logger.Error("Failed", "No Auth Token")
		return false
	}

	succ := addBirthday(task, parsedResponse.AuthToken)

	return succ
}

func addBirthday(task FoodTaskS, auth string) bool {

	headers := client.TLSHeaders{}

	headers.Set("Host", "ordering.api.olo.com")
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("x-device-id", utils.DeviceUUID(true, false, true))
	headers.Set("User-Agent", "Jimmy Johns/5.2.0 (bundle id: com.jimmyjohns.jimmyjohns; build: 20269; iOS: 15.5.0)")
	//cSpell: disable-next-line
	headers.Set("Authorization", "OloKey RQGfKfsvG5kSw0gDN1jCoC967R7GqECd")

	requestData := map[string]any{
		"birthdate":    utils.Birthday("{Y}{M}{D}", false, true),
		"checkrewards": true,
		"checkbalance": true,
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	url := fmt.Sprintf("https://ordering.api.olo.com/v1.1/users/%s/loyaltyschemes/189/provision", auth)

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

	return strings.Contains(string(res), "balance")
}
