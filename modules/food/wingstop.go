package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"strings"
)

func (task FoodTaskS) StartWingstopTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.RandomRealPhoneNumber(),
		Password:  utils.SecurePassword(""),
	}

	res := createWingstopAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func createWingstopAccount(task FoodTaskS) bool {

	headers := client.TLSHeaders{}

	headers.Set("authority", "api.wingstop.com")
	headers.Set("accept", "application/json, text/plain, */*")
	headers.Set("accept-language", "en-US,en;q=0.9")
	headers.Set("clientid", "wingstop")
	headers.Set("content-type", "application/json")
	headers.Set("locale", "en-us")
	headers.Set("nomnom-platform", "web")
	headers.Set("origin", "https://www.wingstop.com")
	headers.Set("referer", "https://www.wingstop.com/")
	headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")

	month, day, _ := utils.BirthdayObject(false, true)

	requestData := map[string]any{
		"emailaddress": task.Account.Email,
		"firstname":    task.Account.FirstName,
		"lastname":     task.Account.LastName,
		"optin":        false,
		"password":     task.Account.Password,
		"nomnom": map[string]any{
			"zip":      utils.RndZip(),
			"dobmonth": month,
			//cSpell:disable
			"dobday":  day,
			"dobyear": "1904",
			"country": "USA",
			//cSpell:enable
		},
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://api.wingstop.com/users/create",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	// println(string(res))

	return strings.Contains(string(res), "authtoken")
}
