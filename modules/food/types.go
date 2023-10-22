package food

import (
	"emerald/client"
	"emerald/utils"
)

type FoodAccount struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type Modules string

const (
	BJS           Modules = "bjs"
	Bojangles     Modules = "bojangles"
	Chilis        Modules = "chilis"
	DelTaco       Modules = "deltaco"
	Dennys        Modules = "dennys"
	DQ            Modules = "dq"
	TacoBell      Modules = "tacobell"
	FirehouseSubs Modules = "firehousesubs"
	IHOP          Modules = "ihop"
	JimmyJohns    Modules = "jimmyjohns"
	KrispyKreme   Modules = "krispykreme"
	McDonalds     Modules = "mcdonalds"
	Panera        Modules = "panera"
	Popeyes       Modules = "popeyes"
	SteakAndShake Modules = "steakandshake"
	Wendys        Modules = "wendys"
	Whataburger   Modules = "whataburger"
	Wingstop      Modules = "wingstop"
	Test          Modules = "test"
)

type FoodTaskS struct {
	Id       string  `json:"id"`
	ModName  Modules `json:"modName"`
	Logger   utils.LoggerS
	Account  FoodAccount
	Proxy    string `json:"proxy"`
	Client   client.HttpClient
	ExtaData ExtaData
	Success  bool
	Error    string
}

type ExtaData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Catchall  string `json:"catchall"`
	Gmail     string `json:"gmail"`
	Use       bool   `json:"use"`
}
