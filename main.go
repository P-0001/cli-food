package main

import (
	"emerald/client"
	"emerald/modules"
	"emerald/modules/food"
	"emerald/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var logger utils.LoggerS = utils.LoggerS{Name: "Home"}
var version = "1.0.1"
var useProxy = false
var debug = false

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func help() {
	modulesArray := enumToStringArray()
	println("Usage: emerald <module>")
	println("Modules: " + strings.Join(modulesArray, " "))
	os.Exit(1)
}

func main() {
	defer timer("main")()

	godotenv.Load(".env")

	if debug {
		modules.Test_()
	}

	if os.Getenv("USE_PROXY") == "1" {
		useProxy = true
		logger.Log("Using proxy")
	}

	modulesArray := enumToStringArray()

	args := os.Args[1:]

	if len(args) < 1 {
		help()
	}

	utils.LogLogo()

	logger.Log("Starting " + version)

	proxyURL := ""

	if useProxy {
		proxyURL = os.Getenv("PROXY_URL")
		re := regexp.MustCompile(`((https?)://\S+)`)
		result := re.MatchString(proxyURL)
		if !result {
			println("Invalid proxy url example: http://host:port or http://user:pass@host:port")
			os.Exit(1)
		}
	}

	modArg := strings.ToLower(args[0])

	modArg = regexp.MustCompile(`[^a-z]`).ReplaceAllString(modArg, "")

	if modArg == "all" {
		RunAllTasks(proxyURL)
		return
	}

	if !utils.Contains(modulesArray, modArg) {
		println("Invalid module " + modArg)
		help()
	}

	var mod food.Modules = food.Modules(modArg)

	if len(args) == 2 {

		var task = getTask(args[1])

		startFoodTask(task, mod, proxyURL)
	} else {
		RunTasks(mod, proxyURL, 1)
	}
}

func startFoodTask(data CLITaskDataS, mod food.Modules, proxy string) {
	task := food.FoodTaskS{Id: data.ID, ModName: mod, Proxy: proxy}
	task.Client = client.GetTLS(task.Proxy, client.TLS_DEFAULT)
	name := fmt.Sprintf("%s/%s", mod, task.Id)
	task.Logger = utils.LoggerS{Name: name}

	task.ExtaData = food.ExtaData{
		FirstName: data.SettingsData.FirstName,
		LastName:  data.SettingsData.LastName,
		Catchall:  data.SettingsData.Catchall,
		Gmail:     data.SettingsData.Gmail,
		Use:       data.SettingsData.Use,
	}

	modules.RunTask(task.ModName, task)

	os.Exit(0)
}

func RunTasks(mod food.Modules, proxy string, num int) {
	var wg sync.WaitGroup

	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()

			hex, _ := utils.GenerateRandomHexString(12)
			numberStr := strconv.Itoa(i + 1)

			id := fmt.Sprintf("%s|%s", numberStr, hex)

			task := food.FoodTaskS{Id: id, ModName: mod, Proxy: proxy}
			task.Client = client.GetTLS(task.Proxy, client.TLS_DEFAULT)
			name := fmt.Sprintf("%s/%s", task.ModName, task.Id)
			task.Logger = utils.LoggerS{Name: name}

			// data from env
			task.ExtaData = getExtaData()

			modules.RunTask(task.ModName, task)
		}(i)
	}

	wg.Wait()

}

func RunAllTasks(proxy string) {

	var wg sync.WaitGroup

	mods := enumToStringArray()

	num := len(mods)

	logger.Log("Running " + strconv.Itoa(num) + " tasks")

	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done() 
			hex, _ := utils.GenerateRandomHexString(12)
			numberStr := strconv.Itoa(i + 1)

			id := fmt.Sprintf("%s|%s", numberStr, hex)

			mod := food.Modules(mods[i])

			task := food.FoodTaskS{Id: id, ModName: mod, Proxy: proxy}
			task.Client = client.GetTLS(task.Proxy, client.TLS_DEFAULT)
			name := fmt.Sprintf("%s/%s", task.ModName, task.Id)
			task.Logger = utils.LoggerS{Name: name}
			task.ExtaData = getExtaData()

			modules.RunTask(task.ModName, task)
		}(i)
	}
	
	wg.Wait()
}

func enumToStringArray() []string {
	result := []string{
		string(food.BJS),
		string(food.Bojangles),
		string(food.Chilis),
		string(food.DelTaco),
		string(food.Dennys),
		string(food.DQ),
		string(food.TacoBell),
		// string(food.FirehouseSubs),
		string(food.IHOP),
		string(food.JimmyJohns),
		string(food.KrispyKreme),
		string(food.McDonalds),
		string(food.Panera),
		string(food.Popeyes),
		string(food.SteakAndShake),
		string(food.Wendys),
		string(food.Whataburger),
		string(food.Wingstop),
	}
	return result
}

func getTask(taskArg string) CLITaskDataS {
	decodedBytes, err := base64.StdEncoding.DecodeString(taskArg)
	if err != nil {
		log.Fatal("Error decoding base64:", err)
	}
	var task CLITaskDataS

	err = json.Unmarshal(decodedBytes, &task)

	if err != nil {
		log.Fatal("Error decoding json:", err)
	}

	return task
}

func getExtaData() food.ExtaData {
	return food.ExtaData{
		FirstName: os.Getenv("FIRST_NAME"),
		LastName:  os.Getenv("LAST_NAME"),
		Catchall:  os.Getenv("CATCHALL"),
		Gmail:     os.Getenv("GMAIL"),
		Use:       os.Getenv("USE_EXTRA_DATA") == "1",
	}
}

type CLITaskDataS struct {
	Type         string       `json:"type"`
	ID           string       `json:"id"`
	UseProxy     bool         `json:"useProxy"`
	SettingsData SettingsData `json:"settingsData"`
}

type SettingsData struct {
	Catchall  string `json:"catchall"`
	Use       bool   `json:"use"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gmail     string `json:"gmail"`
}
