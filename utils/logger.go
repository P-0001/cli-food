package utils

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"

	color "github.com/fatih/color"

	"unicode/utf8"
)

var useColor bool = true

func getTimeStr() string {
	now := time.Now()
	clock := now.Format("03:04:05.000")
	return fmt.Sprintf("[%s]", clock)
}

func Clear() {
	osName := runtime.GOOS
	var cmd *exec.Cmd
	if osName == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type LoggerS struct {
	Name string
}

//this logging stuff is ugly but it works for now...

func colorArrayBlue(arr []any) []any {
	colored := []any{}
	for _, s := range arr {
		if reflect.TypeOf(s).Name() == "string" {
			x := color.BlueString(fmt.Sprintf("%s", s))
			colored = append([]any{x}, colored...)
		} else {
			colored = append([]any{s}, colored...)
		}
	}
	return colored
}

func colorArrayRed(arr []any) []any {
	colored := []any{}
	for _, s := range arr {
		if reflect.TypeOf(s).Name() == "string" {
			x := color.RedString(fmt.Sprintf("%s", s))
			colored = append([]any{x}, colored...)
		} else {
			colored = append([]any{s}, colored...)
		}
	}
	return colored
}

func (l LoggerS) Log(args ...any) {
	var arr []any
	if useColor {
		_time := color.GreenString(getTimeStr())
		name := color.MagentaString(l.Name)
		x := colorArrayBlue(args)
		arr = append([]any{_time, "|", name, "|"}, x...)
		arr2 := append([]any{_time, "|", name, "|"}, args...)
		saveToFile(arr2...)
	} else {
		arr = append([]any{getTimeStr(), l.Name}, args...)
		saveToFile(arr...)
	}
	fmt.Println(arr...)
}

func (l LoggerS) Error(args ...any) {
	var arr []any
	if useColor {
		timeStr := color.GreenString(getTimeStr())
		name := color.MagentaString(l.Name)
		x := colorArrayRed(args)
		arr = []any{timeStr, "|", name, "|"}
		arr = append(arr, x...)
		arr2 := append([]any{timeStr, "|", l.Name, "|"}, args...)
		saveToFile(arr2...)
	} else {
		arr = append([]any{getTimeStr(), l.Name}, args...)
		saveToFile(arr...)
	}

	fmt.Println(arr...)
}

func LogLogo() {
	str := ` 
									                                                  
                          ██████████████████████████████████████                        
                        ██▒▒██    ░░██  ██    ░░██▒▒██    ░░██▓▓██                      
                      ██▒▒  ░░██░░██    ░░██░░██▒▒  ░░██░░██▓▓  ▓▓██                    
                    ██▒▒  ░░  ░░██    ░░  ░░██▒▒  ░░  ░░██▓▓  ▓▓  ▓▓██                  
                    ██████████████████████████████████████████████████                  
                    ██▒▒▒▒▒▒▒▒░░██        ░░██▒▒▒▒▒▒▒▒░░██▓▓▓▓▓▓▓▓▓▓██                  
                      ██▒▒  ░░░░██      ░░░░██▒▒    ░░░░██▓▓  ▓▓▓▓██                    
                        ██▒▒  ░░░░██    ░░░░██▒▒  ░░░░██▓▓  ▓▓▓▓██                      
                          ██▒▒  ░░██    ░░░░██▒▒  ░░░░██▓▓  ▓▓██                        
                            ██▒▒  ░░██  ░░░░██▒▒  ░░██▓▓  ▓▓██                          
                              ██▒▒░░██  ░░░░██▒▒  ░░██▓▓▓▓██                            
        ░░      ░░              ██▒▒░░██  ░░██▒▒░░██▓▓▓▓██                ░░      ░░    
                                  ██▒▒██  ░░██▒▒░░██▓▓██  ░░                            
                                    ██▒▒██░░██▒▒██▓▓██                                  
                                      ████░░██▒▒████                                    
                                    ░░  ██░░██▒▒██        ░░                            
                                          ██████                                                                                                                    
`

	str = color.MagentaString(str)

	fmt.Println(str)
}

func joinArray(arr []any) string {
	str := ""
	for _, s := range arr {
		str += fmt.Sprintf("%v ", s)
	}
	return str
}

func saveToFile(args ...any) {
	// The file path of the text file to append to
	filePath := "log.txt"

	// Text to append to the file
	textToAppend := removeNonUTF8(joinArray(args))

	// Open the file in append mode. If the file doesn't exist, it will be created.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	// Write the text to the file
	_, err = file.WriteString(textToAppend + "\n")

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

func removeNonUTF8(input string) string {
	// Convert the input string to a rune slice to handle Unicode characters
	runeSlice := []rune(input)
	var result []rune

	for _, char := range runeSlice {
		// Check if the rune is a valid UTF-8 encoded Unicode code point
		if utf8.ValidRune(char) {
			result = append(result, char)
		}
	}

	// Convert the result rune slice back to a string
	return string(result)
}
