package modules

import (
	"emerald/modules/food"
	"emerald/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-faker/faker/v4"
)

func needsExtra(mod food.Modules) int {
	switch mod {
	case food.McDonalds:
		return 1
	case food.Popeyes:
		return 2
	case food.Whataburger:
		return 1
	case food.Wingstop:
		return 1
	default:
		return 0
	}
}

func RunTask(mod food.Modules, task food.FoodTaskS) food.FoodTaskS {

	mailType := needsExtra(mod)

	if task.ExtaData.Use {
		if task.ExtaData.FirstName != "" {
			task.Account.FirstName = task.ExtaData.FirstName
		} else {
			task.Account.FirstName = utils.RemoveSpecials(faker.FirstName())
		}

		if task.ExtaData.LastName != "" {
			task.Account.LastName = task.ExtaData.LastName
		} else {
			task.Account.LastName = utils.RemoveSpecials(faker.LastName())
		}
	} else {
		task.Account.FirstName = utils.RemoveSpecials(faker.FirstName())
		task.Account.LastName = utils.RemoveSpecials(faker.LastName())
		// task.Account.Email = utils.RandomEmail(task.Account.FirstName, task.Account.LastName, "")
	}

	// if its no 0 then we need to use a catchall or gmail
	if mailType != 0 {
		if task.ExtaData.Catchall != "" {
			task.Account.Email = utils.RandomEmail(task.Account.FirstName, task.Account.LastName, task.ExtaData.Catchall)
		} else if task.ExtaData.Gmail != "" {
			var jigType utils.JigType

			if utils.RndBool() {
				jigType = "dot"
			} else {
				jigType = "plus"
			}

			task.Account.Email = utils.GmailJig(task.ExtaData.Gmail, jigType)
		} else {
			println("No catchall or gmail provided")
			os.Exit(0)
		}
	} else {
		task.Account.Email = utils.RandomEmail(task.Account.FirstName, task.Account.LastName, "")
	}

	switch mod {

	case food.KrispyKreme:
		task = task.StartKKTask()
		break

	case food.BJS:
		task = task.StartBJSTask()
		break

	case food.SteakAndShake:
		task = task.StartSNSTask()
		break

	case food.Bojangles:
		task = task.StartBojanglesTask()
		break

	case food.TacoBell:
		task = task.StartTacoBellTask()
		break
	case food.Panera:
		task = task.StartPaneraTask()
		break

	case food.Whataburger:
		task = task.StartWhataburgerTask()
		break

	case food.DQ:
		task = task.StartDQTask()
		break

	case food.JimmyJohns:
		task = task.StartJimmyjohnsTask()
		break

	case food.McDonalds:
		task = task.StartMcdonaldsTask()
		break

	case food.Popeyes:
		task = task.StartPopeyesTask()
		break

	case food.Wendys:
		task = task.StartWendysTask()
		break

	case food.Dennys:
		task = task.StartDennysTask()
		break

	case food.IHOP:
		task = task.StartIhopTask()
		break

	case food.Chilis:
		task = task.StartChillisTask()
		break

	case food.DelTaco:
		task = task.StartDeltacoTask()
		break

	case food.Wingstop:
		task = task.StartWingstopTask()
		break

	default:
		task.Logger.Error("Module not found")
		task.Success = false
	}

	/*
		values := reflect.ValueOf(task.Account)
		types := values.Type()
		for i := 0; i < values.NumField(); i++ {
			fmt.Println(types.Field(i).Index[0], types.Field(i).Name, values.Field(i))
		}
	*/

	/*
		task.Logger.Log(fmt.Sprintf("task|%s|%s", task.Id, data))
	*/
	if task.Success && os.Getenv("DEV") == "1" {
		save(task)
	}

	return task
}

func Test_() {

	os.Exit(0)
}

func save(task food.FoodTaskS) {

	filePath := "accounts.txt"

	accountJSON, err := json.Marshal(task.Account)

	if err != nil {
		fmt.Println("Error marshalling account:", err)
		return
	}

	accountString := string(accountJSON)

	textToAppend := fmt.Sprintf("%s|%s", task.ModName, accountString)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(textToAppend + "\n")

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}
