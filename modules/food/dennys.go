package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

func (task FoodTaskS) StartDennysTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createDennysAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

type DennysMakeS struct {
	User        UserS        `json:"user"`
	AccessToken AccessTokenS `json:"access_token"`
	AuthToken   string       `json:"authToken"`
}

type AccessTokenS struct {
	Token           string  `json:"token"`
	SecondsToExpire float64 `json:"seconds_to_expire"`
}

type UserS struct {
	Email               string `json:"email"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Phone               string `json:"phone"`
	EmailMarketingOptIn bool   `json:"email_marketing_opt_in"`
	SMSMarketingOptIn   bool   `json:"sms_marketing_opt_in"`
	TermsAndConditions  bool   `json:"terms_and_conditions"`
	CreatedAt           string `json:"created_at"`
	UserID              string `json:"user_id"`
	AnniversaryYear     int64  `json:"anniversary_year"`
	AnniversaryMonth    int64  `json:"anniversary_month"`
	AnniversaryDay      int64  `json:"anniversary_day"`
}

func createDennysAccount(task FoodTaskS) bool {

	os := utils.RndIosVersion()

	ua := fmt.Sprintf("Mozilla/5.0 (iPhone; CPU iPhone OS %s like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148", os)

	headers := client.TLSHeaders{}

	headers.Set("Host", "nomnom-prod-api.dennys.com")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "capacitor://localhost")
	headers.Set("nomnom-platform", "ios")
	headers.Set("Clientid", "dennys")
	headers.Set("User-Agent", ua)
	headers.Set("Accept-Language", "en-US,en;q=0.9")

	requestData := map[string]any{
		"ignore": nil,
		"user": map[string]any{
			"email":                    task.Account.Email,
			"first_name":               task.Account.FirstName,
			"last_name":                task.Account.LastName,
			"password":                 task.Account.Password,
			"password_confirm":         task.Account.Password,
			"phone":                    task.Account.Phone,
			"terms_and_conditions":     true,
			"sms_marketing_opt_in":     nil,
			"push_notification_opt_in": false,
			"email_marketing_opt_in":   false,
		},
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://nomnom-prod-api.dennys.com/user/create",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	// println(string(res) + "\r\n")

	var response DennysMakeS

	err = json.Unmarshal(res, &response)

	if err != nil {
		task.Logger.Error("Json Unmarshal Failed", err.Error())
		return false
	}

	if response.AccessToken.Token == "" {
		task.Logger.Error("Failed", "No Auth Token")
		return false
	}

	succ := addDennysBirthday(task, response)

	return succ
}
func addDennysBirthday(task FoodTaskS, response DennysMakeS) bool {

	os := utils.RndIosVersion()

	ua := fmt.Sprintf("Mozilla/5.0 (iPhone; CPU iPhone OS %s like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148", os)

	headers := client.TLSHeaders{}
	headers.Set("User-Agent", ua)
	headers.Set("Host", "nomnom-prod-api.dennys.com")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "capacitor://localhost")
	headers.Set("nomnom-platform", "ios")
	headers.Set("Clientid", "dennys")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", response.AccessToken.Token))

	d, m, y := utils.BirthdayObject(false, false)

	requestData := map[string]any{
		"user": map[string]any{
			"user_id":                  response.User.UserID,
			"first_name":               task.Account.FirstName,
			"last_name":                task.Account.LastName,
			"email":                    task.Account.Email,
			"email_marketing_opt_in":   false,
			"sms_marketing_opt_in":     nil,
			"phone":                    task.Account.Phone,
			"favorite_location_ids":    "9525",
			"terms_and_conditions":     true,
			"birth_year":               utils.ToInt(y),
			"birth_month":              utils.ToInt(m),
			"birth_day":                utils.ToInt(d),
			"anniversary_year":         nil,
			"anniversary_month":        nil,
			"anniversary_day":          nil,
			"push_notification_opt_in": false,
		},
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://nomnom-prod-api.dennys.com/user",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), fmt.Sprintf(`"birth_year":%s`, y))
}
