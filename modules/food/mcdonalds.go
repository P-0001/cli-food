package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"strings"
)

const clientSecret = "Ym4rVyqpqNpCpmrdPGJatRrBMHhJgr26"
const clientId = "8cGckR5wPgQnFBc9deVhJ2vT94WhMBRL"
const basicAuth = "Basic OGNHY2tSNXdQZ1FuRkJjOWRlVmhKMnZUOTRXaE1CUkw6WW00clZ5cXBxTnBDcG1yZFBHSmF0UnJCTUhoSmdyMjY="
const mcVersion = "7.11.2"
const mcSDKVersion = "24.0.18"

func (task FoodTaskS) StartMcdonaldsTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createMcdonaldsAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account Created")
	} else {
		task.Logger.Error("Account Creation Failed")
	}

	return task
}

func createMcdonaldsAccount(task FoodTaskS) bool {

	deviceId := utils.DeviceUUID(true, false, true)

	succ, auth := genAuthToken(task)

	if !succ {
		task.Logger.Error("Failed to generate auth token")
		return false
	}

	succ = postEmail(task, auth, deviceId)

	if !succ {
		task.Logger.Error("Failed to post email")
		return false
	}

	headers := client.TLSHeaders{}
	headers.Set("Host", "us-prod.api.mcd.com")
	headers.Set("mcd-sourceapp", "GMA")
	headers.Set("cache-control", "true")
	headers.Set("user-agent", getUA())
	headers.Set("mcd-clientsecret", clientSecret)
	headers.Set("mcd-uuid", utils.DeviceUUID(true, false, true))
	headers.Set("authorization", fmt.Sprintf("Bearer %s", auth))
	headers.Set("mcd-clientid", clientId)
	headers.Set("accept-language", "en-US")
	headers.Set("x-acf-sensor-data", "1,i,WbZByZAQLF/Rlqx2kzdp4mc0RDyluLiYCxZIttB4NT+IGpEOk+4YnJGp1Hrn+pSPDPFSdOMQbJNNnTHcu6Vm30cNeUmV8BPatby/jymVpr3aB3ou+2Eh7LJ7UR4E8b1FRBy6IUyykawKSMYZsDnIYiLGxYSyQTrINLSkE4ZIMYk=,XYhw4QiuybIj74cxhRRBKwqOfMo59PV0b84AZ/nzhxdXVG4T/V9YFGRgPM0fUFgBzV9+JGwj55pTxpb+T4RTrNVu9g4xDXz93ncxVhTNi1ldomMdahmq57Zqi+MQFZ1QE5n/oz8tN8E9K05UUPCc0V0TXRjNTfAL4NyqQrDxM6E=$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$13,6,32")

	requestData := makeBody(deviceId, task.Account.Email, task.Account.FirstName, task.Account.LastName, utils.RndZip())

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://us-prod.api.mcd.com/exp/v1/customer/registration",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	strRes := string(res)

	return strings.Contains(strRes, "successful")
}

func genAuthToken(task FoodTaskS) (bool, string) {

	headers := client.TLSHeaders{}

	headers.Set("Host", "us-prod.api.mcd.com")
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Mcd-Sourceapp", "GMA")
	headers.Set("Accept", "application/json")
	headers.Set("Authorization", basicAuth)
	headers.Set("Accept-Charset", "utf-8")
	headers.Set("Mcd-Clientid", clientId)
	headers.Set("Accept-Language", "en-US")
	headers.Set("Cache-Control", "true")
	headers.Set("Mcd-Uuid", utils.DeviceUUID(true, false, true)) // Assuming DeviceId() returns a UUID string
	headers.Set("User-Agent", getUA())                           // Assuming getUA() returns the User-Agent string
	headers.Set("Mcd-Marketid", "US")
	headers.Set("Mcd-Clientsecret", clientSecret)

	body := utils.ToForm(map[any]any{
		"grant_type": "client_credentials",
	})

	reader := bytes.NewReader([]byte(body))

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://us-prod.api.mcd.com/v1/security/auth/token",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false, ""
	}

	type Response struct {
		Token string `json:"token"`
		// Add other fields from the JSON response if needed
	}

	var jsonResponse struct {
		Response Response `json:"response"`
	}

	err = json.Unmarshal(res, &jsonResponse)

	if err != nil {
		task.Logger.Error("json.Unmarshal Failed", err.Error())
		return false, ""
	}

	return true, jsonResponse.Response.Token
}

func postEmail(task FoodTaskS, auth string, deviceId string) bool {
	headers := client.TLSHeaders{}

	headers.Set("Host", "us-prod.api.mcd.com")
	headers.Set("Mcd-Sourceapp", "GMA")
	headers.Set("Cache-Control", "true")
	headers.Set("User-Agent", getUA()) // Assuming getUA() returns the User-Agent string
	headers.Set("Mcd-Clientsecret", clientSecret)
	headers.Set("Mcd-Uuid", utils.DeviceUUID(true, false, true)) // Assuming DeviceId() returns a UUID string
	headers.Set("Authorization", "Bearer "+auth)                 // Assuming authToken is a token string
	headers.Set("Mcd-Clientid", clientId)
	headers.Set("X-Acf-Sensor-Data", "")
	headers.Set("Accept-Language", "en-US")
	headers.Set("Accept-Charset", "utf-8")
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	headers.Set("Mcd-Marketid", "US")

	requestData := map[string]any{
		"deviceId":           deviceId,
		"customerIdentifier": task.Account.Email,
		"registrationType":   "traditional",
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	reader := bytes.NewReader(body)

	_, err = client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://us-prod.api.mcd.com/exp/v1/customer/identity/email",
		Headers:          headers,
		Body:             reader,
		Client:           task.Client,
		ExpectedResponse: 404,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	return true
}

func getUA() string {
	iosV := utils.RndIosVersion()

	return fmt.Sprintf("MCDSDK/%s (iPhone; %s; en-US) GMA/%s", mcSDKVersion, iosV, mcVersion)
}

func makeBody(deviceID, email, firstName, lastName, zip string) any {
	type AccountData struct {
		email     string
		firstName string
		lastName  string
	}

	type PreferencesDetails struct {
		Email     string `json:"email,omitempty"`
		MobileApp string `json:"mobileApp,omitempty"`
		Enabled   string `json:"enabled,omitempty"`
	}

	type PreferencesItem struct {
		Details      PreferencesDetails `json:"details"`
		PreferenceID int                `json:"preferenceId"`
	}

	type SubscriptionsItem struct {
		SubscriptionID string `json:"subscriptionId"`
		OptInStatus    string `json:"optInStatus"`
	}

	type MakeBodyRequest struct {
		Policies struct {
			AcceptancePolicies map[string]bool `json:"acceptancePolicies"`
		} `json:"policies"`
		OptInForMarketing bool              `json:"optInForMarketing"`
		EmailAddress      string            `json:"emailAddress"`
		Preferences       []PreferencesItem `json:"preferences"`
		Address           struct {
			ZipCode string `json:"zipCode"`
			Country string `json:"country"`
		} `json:"address"`
		Device struct {
			DeviceID     string `json:"deviceId"`
			OS           string `json:"os"`
			OSVersion    string `json:"osVersion"`
			DeviceIDType string `json:"deviceIdType"`
			IsActive     string `json:"isActive"`
			Timezone     string `json:"timezone"`
		} `json:"device"`
		Credentials struct {
			Type          string `json:"type"`
			LoginUsername string `json:"loginUsername"`
			SendMagicLink bool   `json:"sendMagicLink"`
		} `json:"credentials"`
		Audit struct {
			RegistrationChannel string `json:"registrationChannel"`
		} `json:"audit"`
		Application   string              `json:"application"`
		FirstName     string              `json:"firstName"`
		Subscriptions []SubscriptionsItem `json:"subscriptions"`
		LastName      string              `json:"lastName"`
	}

	return MakeBodyRequest{
		Policies: struct {
			AcceptancePolicies map[string]bool `json:"acceptancePolicies"`
		}{
			AcceptancePolicies: map[string]bool{
				"1": true,
				"6": false,
				"4": true,
				"5": false,
			},
		},
		OptInForMarketing: false,
		EmailAddress:      email,
		Preferences: []PreferencesItem{
			{
				Details: PreferencesDetails{
					Email:     "en-US",
					MobileApp: "en-US",
				},
				PreferenceID: 1,
			},
			{
				Details: PreferencesDetails{
					Email:     "0",
					MobileApp: "Y",
				},
				PreferenceID: 2,
			},
			{
				Details: PreferencesDetails{
					Email:     "Y",
					MobileApp: "Y",
				},
				PreferenceID: 3,
			},
			// Add other preference items here
		},
		Address: struct {
			ZipCode string `json:"zipCode"`
			Country string `json:"country"`
		}{
			ZipCode: zip,
			Country: "US",
		},
		Device: struct {
			DeviceID     string `json:"deviceId"`
			OS           string `json:"os"`
			OSVersion    string `json:"osVersion"`
			DeviceIDType string `json:"deviceIdType"`
			IsActive     string `json:"isActive"`
			Timezone     string `json:"timezone"`
		}{
			DeviceID:     deviceID,
			OS:           "ios",
			OSVersion:    "15.5",
			DeviceIDType: "IDFV",
			IsActive:     "Y",
			Timezone:     "America/Los_Angeles",
		},
		Credentials: struct {
			Type          string `json:"type"`
			LoginUsername string `json:"loginUsername"`
			SendMagicLink bool   `json:"sendMagicLink"`
		}{
			Type:          "email",
			LoginUsername: email,
			SendMagicLink: true,
		},
		Audit: struct {
			RegistrationChannel string `json:"registrationChannel"`
		}{
			RegistrationChannel: "M",
		},
		Application: "gma",
		FirstName:   firstName,
		Subscriptions: []SubscriptionsItem{
			{
				SubscriptionID: "1",
				OptInStatus:    "Y",
			},
			{
				SubscriptionID: "2",
				OptInStatus:    "Y",
			},
			{
				SubscriptionID: "3",
				OptInStatus:    "Y",
			},
			// Add other subscriptions here
		},
		LastName: lastName,
	}
}
