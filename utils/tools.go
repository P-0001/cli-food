package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	"encoding/base64"
)

func RndInt(_max int) int {
	setSeed()
	return rand.Intn(_max)
}

func RndRangeInt(_min int, _max int) int {
	setSeed()
	return rand.Intn(_max-_min+1) + _min
}

func InsertStringAtIndex(index int, original string, insert string) string {
	return original[:index] + insert + original[index:]
}

func RandomCharFromString(chars string) string {
	return string(chars[RndInt(len(chars))])
}

func RandomDomain() string {
	domains := [8]string{
		"gmail.com",
		"outlook.com",
		"yahoo.com",
		"hotmail.com",
		"aol.com",
		"comcast.net",
		"att.net",
		"icloud.com",
	}
	return domains[RndInt(8)]
}

func RandomRealPhoneNumber() string {
	phone := RandomAreaCode()

	for i := 0; i < 7; i++ {
		phone += RandomCharFromString("0123456789")
	}

	if len(phone) != 10 {
		panic("Invalid phone number")
	}

	return phone
}

func FormatPhoneNumber(q string, format string) string {
	raw := regexp.MustCompile(`[^0-9]`).ReplaceAllString(q, "")
	formatted := make([]string, len(format))
	index := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '#' {
			formatted[i] = string(raw[index])
			index++
		} else {
			formatted[i] = string(format[i])
		}
	}

	return strings.Join(formatted, "")
}

var specials = "!@"

func SecurePassword(_specials string) string {
	password := ""

	password += RandomString("abcdefghijklmnopqrstuvwxyz", 4)

	password += RandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4)

	password += RandomCharFromString("0123456789")

	if len(specials) < 0 {
		password += RandomCharFromString(_specials)
	} else {
		password += RandomCharFromString(specials)
	}

	return ShuffleString(password)
}

func RandomString(chars string, length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s += RandomCharFromString(chars)
	}
	return s
}

func ShuffleString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := n - 1; i > 0; i-- {
		j := RndInt(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func RandomEmail(firstName, lastName, domain string) string {
	if firstName == "" {
		firstName = faker.FirstName()
	}
	if lastName == "" {
		lastName = faker.LastName()
	}
	if domain == "" {
		domain = RandomDomain()
	}

	suffix := RandomString("abcdefghijklmnopqrstuvwxyz1234567890", RndRangeInt(3, 8))

	base := firstName + lastName + suffix

	raw := regexp.MustCompile(`[^0-9a-zA-Z]`).ReplaceAllString(base, "")

	return raw + "@" + domain
}

func RndIosVersion() string {
	major := RndRangeInt(12, 16)
	_min := RndInt(10)

	return fmt.Sprintf("%d.%d", major, _min)
}

func RndIosVersion3() string {
	major := RndRangeInt(12, 16)
	_min := RndInt(10)
	pat := RndInt(10)

	return fmt.Sprintf("%d.%d.%d", major, _min, pat)
}

func GenerateRandomHexString(length int) (string, error) {
	byteLength := length / 2

	randomBytes := make([]byte, byteLength)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	hexString := hex.EncodeToString(randomBytes)

	if len(hexString) > length {
		hexString = hexString[:length]
	}

	return hexString, nil
}

func DeviceUUID(slashes bool, lower bool, upper bool) string {
	id := uuid.NewString()

	if !slashes {
		id = regexp.MustCompile(`-`).ReplaceAllString(id, "")
	}

	if lower {
		id = strings.ToLower(id)
	}

	if upper {
		id = strings.ToUpper(id)
	}

	return id
}

func Contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

type JigType string

// Define constants representing the allowed values
const (
	Plus JigType = "plus"
	Dot  JigType = "dot"
)

func GmailJig(gmail string, jigType JigType) string {

	arr := strings.Split(gmail, "@")

	meat := arr[0]

	newString := meat

	if jigType == Plus {
		newString = meat + "+" + RandomString("abcdefghijklmnopqrstuvwxyz1234567890", RndRangeInt(3, 8))
	} else {
		length := len(meat)
		amount := RndInt(2) + 1
		sent := 0
		for i := 0; i < length; i++ {
			index := RndInt(length)
			if index == 0 || index == length-1 {
				continue
			}
			if meat[index-1] == '.' || meat[index+1] == '.' {
				continue
			}
			if sent < amount {
				newString = InsertStringAtIndex(index, newString, ".")
				sent++
				length++
			}
		}
	}

	return newString + "@gmail.com"
}

func ToInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return number
}

func RemoveSpecials(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(str, "")
}

func RndBool() bool {
	return RndInt(2) == 0
}

func ToString(v any) string {
	return fmt.Sprintf("%v", v)
}

func DecodeB64(base64Str string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)

	if err != nil {
		return nil, err
	}

	return decodedBytes, nil
}

func GetCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}

var seedSet = false

func setSeed() {
	if !seedSet {
		rand.Seed(time.Now().UnixNano())
		seedSet = true
	}
}
