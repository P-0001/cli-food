package food

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"emerald/client"
	"emerald/utils"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func (task FoodTaskS) StartSNSTask() FoodTaskS {
	task.Logger.Log("Starting Account Creation")

	task.Account = FoodAccount{
		FirstName: task.Account.FirstName,
		LastName:  task.Account.LastName,
		Email:     task.Account.Email,
		Phone:     utils.FormatPhoneNumber(utils.RandomRealPhoneNumber(), "(###) ###-####"),
		Password:  utils.SecurePassword(""),
	}

	res := createSNSAccount(task)

	task.Success = res

	if res {
		task.Logger.Log("Account created")
	} else {
		task.Logger.Error("Account creation failed")
	}

	return task
}

type userDataS struct {
	Client string `json:"client"`
	User   user   `json:"user"`
}

type user struct {
	Birthday             string `json:"birthday"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	FavouriteLocationIDs string `json:"favourite_location_ids"`
	Phone                string `json:"phone"`
	TermsAndConditions   bool   `json:"terms_and_conditions"`
}

func createSNSAccount(task FoodTaskS) bool {
	requestData := userDataS{
		Client: "89f14e9b57c95083cd91b307c6e7c4682e56ab12765cb2260700062bbcb991e3",
		User: user{
			Birthday:             utils.Birthday("{Y}-{M}-{D}", false, true),
			Email:                task.Account.Email,
			FirstName:            task.Account.FirstName,
			LastName:             task.Account.LastName,
			Password:             task.Account.Password,
			PasswordConfirmation: task.Account.Password,
			FavouriteLocationIDs: "345555",
			Phone:                task.Account.Phone,
			TermsAndConditions:   true,
		},
	}

	body, err := json.Marshal(requestData)

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	payload := "/api2/mobile/users" + string(body)

	digest := hash([]byte(payload))

	headers := client.TLSHeaders{
		"Host":                 []string{"mobileandroidapi.punchh.com"},
		"x-pch-digest":         []string{digest},
		"user-agent":           []string{"com.zipscene.mobile.sns/4.2.1/1422(Android;Samsung;Samsung Galaxy S10;10;xxxhdpi)"},
		"punchh-app-device-id": []string{uuid.NewString()},
		"content-type":         []string{"application/json; charset=UTF-8"},
	}

	res, err := client.TlsRequest(client.TLSParams{
		Method:           "POST",
		Url:              "https://mobileandroidapi.punchh.com/api2/mobile/users",
		Headers:          headers,
		Body:             bytes.NewReader(body),
		Client:           task.Client,
		ExpectedResponse: 200,
	})

	if err != nil {
		task.Logger.Error("Failed", err.Error())
		return false
	}

	return strings.Contains(string(res), "access_token")
}

func hash(message []byte) string {
	key := []byte("7e9c5606a39142c0eb05f86f4f42f3801f22b5d9157d9a92d7242e21d5cc55ab") // Replace this with your actual secret key

	// Create a new HMAC-SHA256 hasher
	hasher := hmac.New(sha256.New, key)

	// Write the message to the hasher
	hasher.Write(message)

	// Get the computed HMAC-SHA256 hash as a byte slice
	hash := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashHex := hex.EncodeToString(hash)

	return hashHex
}
