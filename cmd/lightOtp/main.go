package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codemicro/lightOtp/internal/commands"
	"github.com/codemicro/lightOtp/internal/helpers"
	"github.com/codemicro/lightOtp/internal/persist"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	VERSION string = "0.0.0"
)

func main() {
	col := color.New(color.FgCyan, color.Bold)
	_, _ = col.Print(" _ _       _     _     ___ _         \n| (_) __ _| |__ | |_  /___\\ |_ _ __  \n| | |/ _` |" +
		" '_ \\| __|//  // __| '_ \\ \n| | | (_| | | | | |_/ \\_//| |_| |_) |\n|_|_|\\__, |_| |_|\\__\\___/  \\__| ._" +
		"_/ \n     |___/                    |_| ")
	fmt.Println("v" + VERSION)
	fmt.Println()

	// Load Settings (create new with defaults if does not exist)

	settingsFileContent, settingsFileLocation, err := helpers.OpenConfigFile("Settings.json")
	helpers.CheckErr(err)

	if len(settingsFileContent) == 0 {
		helpers.PrintInfoLn("Settings file is empty or missing - creating new with default values at " +
			settingsFileLocation)

		// Create new Settings

		persist.Settings, err = helpers.NewSettings()
		helpers.CheckErr(err)

		fCont, _ := json.Marshal(&persist.Settings)
		err = ioutil.WriteFile(settingsFileLocation, fCont, 0644)
		helpers.CheckErr(err)

		settingsFileContent, _, err = helpers.OpenConfigFile("settings.json")
		helpers.CheckErr(err)

	} else {
		err = json.Unmarshal([]byte(settingsFileContent), &persist.Settings)
		helpers.QuitWithMessageIfErr(err, "Unable to parse JSON in the settings file. Quitting.")
	}

	if clipboard.Unsupported {
		helpers.PrintInfoLn("Writing to the clipboard is unavailable. See [INSERT LINK HERE].")
	}

	// Load codes

	rawCodesJson, err := ioutil.ReadFile(persist.Settings.CodesLocation)
	if err != nil { // Assuming it means CodesLocation does not exist
		helpers.PrintInfoLn("Cannot find codes file - creating new from scratch")

		file, err := os.OpenFile(persist.Settings.CodesLocation, os.O_RDONLY|os.O_CREATE, 0666)
		helpers.QuitWithMessageIfErr(err, "Unable to create codes file. Quitting.")
		defer file.Close()

		_, _ = file.Write([]byte("[]"))

		rawCodesJson = []byte("[]")

	}

	err = json.Unmarshal(rawCodesJson, &persist.Codes)
	helpers.QuitWithMessageIfErr(err, "Unable to parse JSON in the codes file. Quitting.")

	// Main program loop

	scanner := bufio.NewScanner(os.Stdin)
	for {
		_, _ = color.New(color.FgCyan).Print("> ")
		scanner.Scan()

		text := scanner.Text()
		splitText := strings.Split(text, " ")

		switch splitText[0] {
		case "help":
			commands.Help()
		case "list":
			commands.ListProviders()
		case "code":
			if len(splitText) < 2 {
				helpers.PrintErrLn("Not enough arguments")
			} else {
				i, err := strconv.ParseInt(splitText[1], 10, 32)
				if err != nil {
					helpers.PrintErrLn(splitText[1] + ": invalid number")
				} else {
					commands.GenerateCode(int32(i - 1))
				}
			}
		case "add":
			commands.AddProvider()
		case "del":
			if len(splitText) < 2 {
				helpers.PrintErrLn("Not enough arguments")
			} else {
				i, err := strconv.ParseInt(splitText[1], 10, 32)
				if err != nil {
					helpers.PrintErrLn(splitText[1] + ": invalid number")
				} else {
					commands.RemoveProvider(int32(i - 1))
				}
			}
		case "exit":
			fmt.Println("Bye o/")
			os.Exit(0)
		default:
			helpers.PrintErrLn(text + ": unknown command. Try running 'help'.")
		}

		fmt.Println()
	}

}
