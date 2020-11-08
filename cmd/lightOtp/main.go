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
	"unicode/utf8"
)

const (
	VERSION string = "0.0.0"
)

var (
	scanner = bufio.NewScanner(os.Stdin)
)

func main() {
	col := color.New(color.FgCyan, color.Bold)
	_, _ = col.Print(" _ _       _     _     ___ _         \n| (_) __ _| |__ | |_  /___\\ |_ _ __  \n| | |/ _` |" +
		" '_ \\| __|//  // __| '_ \\ \n| | | (_| | | | | |_/ \\_//| |_| |_) |\n|_|_|\\__, |_| |_|\\__\\___/  \\__| ._" +
		"_/ \n     |___/                    |_| ")
	fmt.Println("v" + VERSION)
	fmt.Println()

	// Load Settings (create new with defaults if does not exist)

	settingsFileContent, settingsFileLocation, err := helpers.OpenConfigFile("settings.json")
	helpers.CheckErr(err)

	if utf8.RuneCountInString(settingsFileContent) == 0 {
		helpers.PrintInfoLn("Settings file is empty or missing - creating new with default values at " +
			settingsFileLocation)

		// Create new Settings

		persist.Settings = helpers.NewSettings() // Create new object with defaults

		fCont, _ := json.Marshal(&persist.Settings)
		err = ioutil.WriteFile(settingsFileLocation, fCont, 0644)
		helpers.QuitWithMessageIfErr(err, "Unable to save settings file")

	} else {
		err = json.Unmarshal([]byte(settingsFileContent), &persist.Settings)
		helpers.QuitWithMessageIfErr(err, "Unable to parse JSON in the settings file. Quitting.")
	}

	if clipboard.Unsupported {
		helpers.PrintInfoLn("Writing to the clipboard is unavailable. Check the README.")
	}

	// Load codes

	if _, err := os.Stat(persist.Settings.CodesLocation); os.IsNotExist(err) {
		helpers.PrintInfoLn("Cannot find codes file - creating new from scratch")

		for {
			_, _ = color.New(color.FgCyan).Print("Please set a master password > ")
			firstPassword := helpers.CollectCensoredInput()
			_, _ = color.New(color.FgCyan).Print("Repeat > ")
			if firstPassword != helpers.CollectCensoredInput() {
				helpers.PrintErrLn("Passwords do not match. Try again.")
				fmt.Println()
			} else {
				persist.MasterPassword = firstPassword
				break
			}
		}

		err := helpers.UpdateCodes()
		helpers.QuitWithMessageIfErr(err, "Unable to create codes file. Quitting.")

		helpers.PrintInfoLn("File creation successful!")

	} else {
		// Collect master password
		_, _ = color.New(color.FgCyan).Print("Master password > ")
		persist.MasterPassword = helpers.CollectCensoredInput()
	}

	err = helpers.LoadCodes()
	helpers.QuitWithMessageIfErr(err, "Unable to load codes file. Quitting. (Is the password incorrect?)")

	helpers.PrintInfoLn("Password ok")

	fmt.Println()

	// Main program loop

	for {
		_, _ = color.New(color.FgCyan).Print("> ")
		scanner.Scan()

		text := scanner.Text()
		splitText := strings.Split(text, " ")

		switch splitText[0] {
		case "help":
			commands.Help()
		case "list", "ls":
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
		case "change":
			commands.ChangePassword()
		default:
			helpers.PrintErrLn(text + ": unknown command. Try running 'help'.")
		}

		fmt.Println()
	}

}
