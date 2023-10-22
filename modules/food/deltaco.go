package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"strings"
)

func (task FoodTaskS) StartDeltacoTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "### ### ####"),
		Password:  utils.SecurePassword(""),
	}

	// println("phone", task.Account.Phone)

	res := createDeltacoAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createDeltacoAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}
	headers.Set("Host", "cust1184.cheetahedp.com")
	headers.Set("Accept", "application/vnd.stellar-v1+json")
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("User-Agent", "DelTaco/3.7.4 iOS/15.5/iPhone")
	headers.Set("Accept-Language", "en-US,en;q=0.9")

	requestData := map[any]any{
		"accept_terms_and_conditions": "1",
		"birthdate":                   utils.Birthday("1998-{M}-{D}", false, true),
		"client_id":                   "20c490787f452826151dde68dfdc5ed32112a24826399f8bc4d6e62c70ba74ff",
		"client_secret":               "409975196a96346f6c3801233396f499ed06aca2494be0ddf9c9a827002bd444",
		"email":                       task.Account.Email,
		"favorite_location":           "886",
		"favorite_vendor_id":          "25515",
		"first_name":                  task.Account.FirstName,
		"last_name":                   task.Account.LastName,
		"mobile_phone":                task.Account.Phone,
		"password":                    task.Account.Password,
		"password_confirmation":       task.Account.Password,
		"receive_mail_offers":         "0",
		"receive_mobile_app_offers":   "1",
		"receive_sms_offers":          "0",
	}

	body := []byte(utils.ToForm(requestData))

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://cust1184.cheetahedp.com//api/sign_up",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 201,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	if !strings.Contains(string(res), `"success":true`) {
		task.Logger.Error("Failed", string(res))
		return false
	}

	// println("make", string(res))

	// succ := getToken(task)

	return true
}

func getToken(task FoodTaskS) bool {
	headers := client.TLSHeaders{}

	requestData := map[any]any{
		"client_id":     "20c490787f452826151dde68dfdc5ed32112a24826399f8bc4d6e62c70ba74ff",
		"client_secret": "409975196a96346f6c3801233396f499ed06aca2494be0ddf9c9a827002bd444",
		"email":         task.Account.Email,
		"grant_type":    "password",
		"password":      task.Account.Password,
	}

	body := []byte(utils.ToForm(requestData))

	reader := bytes.NewReader(body)

	_, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://cust1184.cheetahedp.com//oauth/token",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 401,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	// println("token", string(res))

	return true
}
