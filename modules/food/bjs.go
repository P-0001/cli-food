package food

import (
	"bytes"
	"emerald/client"
	"emerald/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"

	"github.com/xxtea/xxtea-go/xxtea"
)

var versionBJS = "3.9.3"
var keyBjs = "705F43E7449B4C0EA26FDEE698D7651D"

func (task FoodTaskS) StartBJSTask() FoodTaskS {

	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createBJSAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

func encryptToString(str string) string {
	enc := xxtea.Encrypt([]byte(str), []byte(keyBjs))
	return base64.StdEncoding.EncodeToString(enc)
}

func decryptToString(str string) string {
	dec, _ := base64.StdEncoding.DecodeString(str)
	decrypted := xxtea.Decrypt(dec, []byte(keyBjs))
	return string(decrypted)
}

func genAuth(task FoodTaskS, androidId string) (bool, string) {
	ts := utils.FormatDate("{m}{d}{y}{h}{i}{s}", time.Now(), false)

	str := fmt.Sprintf("%s|%s|GetAuthToken|%s", androidId, versionBJS, ts)

	authorization := encryptToString(str)

	headers := client.TLSHeaders{}
	headers.Add("Host", "services1.bjsrestaurants.com")
	headers.Add("securitytoken", "")
	headers.Add("customerid", "0")
	headers.Add("loyaltyid", "0")
	headers.Add("deviceid", androidId)
	headers.Add("authorization", authorization)
	headers.Add("content-type", "text/plain; charset=ISO-8859-1")
	headers.Add("user-agent", "okhttp/5.0.0-alpha.2")

	body := fmt.Sprintf(`{"OSType":"Android","AppVersion":"%s"}`, versionBJS)

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://services1.bjsrestaurants.com/BJSS/Account.svc/GetAuthToken",
		Headers:          headers,
		Body:             bytes.NewReader([]byte(body)),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed to get auth token", err.Error())
		return false, ""
	}

	type responseData struct {
		GetAuthTokenResult struct {
			Data []struct {
				AuthToken string `json:"AuthToken"`
			} `json:"Data"`
		} `json:"GetAuthTokenResult"`
	}

	var resData responseData

	err = json.Unmarshal([]byte(res), &resData)

	if err != nil {
		task.Logger.Error("Failed to unmarshal auth token")
		return false, ""
	}

	var rawAuthToken string

	if len(resData.GetAuthTokenResult.Data) > 0 {
		rawAuthToken = resData.GetAuthTokenResult.Data[0].AuthToken
	} else {
		task.Logger.Error("Failed to get auth token")
		return false, ""
	}

	authToken := decryptToString(rawAuthToken)

	return true, authToken

}

func createBJSAccount(task FoodTaskS) bool {

	androidId, err := utils.GenerateRandomHexString(16)

	if err != nil {
		task.Logger.Error("Failed to generate random hex string")
		return false
	}

	succ, auth := genAuth(task, androidId)

	if !succ {
		task.Logger.Error("Failed to generate auth token")
		return false
	}

	headers := client.TLSHeaders{}

	headers.Add("Host", "services1.bjsrestaurants.com")
	headers.Add("securitytoken", "")
	headers.Add("customerid", "0")
	headers.Add("loyaltyid", "0")
	headers.Add("deviceid", androidId)
	headers.Add("authorization", auth)
	headers.Add("content-type", "text/plain; charset=ISO-8859-1")
	headers.Add("user-agent", "okhttp/5.0.0-alpha.2")

	address := faker.GetRealAddress()

	requestData := map[string]any{
		"LoyaltyId":             "",
		"Password":              task.Account.Password,
		"FirstName":             task.Account.FirstName,
		"LastName":              task.Account.LastName,
		"BirthDate":             utils.Birthday("{M}/{D}/{Y}", false, true),
		"Email":                 task.Account.Email,
		"Phone":                 task.Account.Phone,
		"SiteId":                554,
		"Over18":                "Y",
		"Allergies":             "",
		"FavoriteFood":          "",
		"FavoriteBeverage":      "",
		"KidsUnder12":           0,
		"PointBalance":          0,
		"SMSOptIn":              "False",
		"EmailOptIn":            "False",
		"EmailLoyaltyOptIn":     "True",
		"SpecialRequest":        "",
		"ZipCode":               utils.RndZip(),
		"BCAddress1":            address.Address,
		"BCAddress2":            "",
		"BCCity":                address.City,
		"BCState":               address.State,
		"BCPickupLocation":      "",
		"BCPaymentOption":       "",
		"CCNumber":              "",
		"CCExpDate":             "",
		"CCZipCode":             "",
		"CCCVVNumber":           "",
		"BCEmailOptInMarketing": "",
		"BCEmailOptInBeerClub":  "",
		"DeviceType":            "phone",
		"DeviceName":            "Android",
		"DeviceUUID":            androidId,
		"DeviceOSVersion":       fmt.Sprintf("REL %d", utils.RndRangeInt(5, 10)),
		"AppVersion":            versionBJS,
		"CCMobilePayment":       "False",
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://services1.bjsrestaurants.com/BJSS/Account.svc/RegisterP3",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("json.Marshal Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), `"Status":"1"`)
}
