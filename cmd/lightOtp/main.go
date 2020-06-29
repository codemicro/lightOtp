package main

import (
	"encoding/json"
	"fmt"
	"github.com/codemicro/lightOtp/internal/helpers"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

const (
	VERSION string = "0.0.0"
)

var (
	settings models.Settings
)

func main() {
	col := color.New(color.FgBlue, color.Bold)
	_, _ = col.Print(" _ _       _     _     ___ _         \n| (_) __ _| |__ | |_  /___\\ |_ _ __  \n| | |/ _` |" +
		" '_ \\| __|//  // __| '_ \\ \n| | | (_| | | | | |_/ \\_//| |_| |_) |\n|_|_|\\__, |_| |_|\\__\\___/  \\__| ._" +
		"_/ \n     |___/                    |_| ")
	fmt.Println("v" + VERSION)

	// Load settings (create new with defaults if does not exist)

	settingsFileContent, settingsFileLocation, err := helpers.OpenConfigFile("settings.json")
	helpers.CheckErr(err)

	if len(settingsFileContent) == 0 {
		helpers.PrintInfoLn("Settings file is empty or missing - creating new with default values at " +
			settingsFileLocation)

		// Create new settings

		settings, err = models.NewSettings()
		helpers.CheckErr(err)

		fCont, _ := json.Marshal(&settings)
		err = ioutil.WriteFile(settingsFileLocation, fCont, 0644)
		helpers.CheckErr(err)

		settingsFileContent, _, err = helpers.OpenConfigFile("settings.json")
		helpers.CheckErr(err)

	} else {
		err = json.Unmarshal([]byte(settingsFileContent), &settings)
		helpers.QuitWitMessageIfErr(err, "Unable to parse JSON in the settings file. Quitting.")
	}

	// Load codes

	rawCodesJson, err := ioutil.ReadFile(settings.CodesLocation)
	if err != nil { // Assuming it means CodesLocation does not exist
		helpers.PrintInfoLn("Cannot find codes file - creating new from scratch")

		file, err := os.OpenFile(settings.CodesLocation, os.O_RDONLY|os.O_CREATE, 0666)
		helpers.QuitWitMessageIfErr(err, "Unable to create codes file. Quitting.")
		defer file.Close()

		_, _ = file.Write([]byte("[]"))

		rawCodesJson = []byte("[]")

	}

	var codes []models.TOTPCode
	err = json.Unmarshal(rawCodesJson, &codes)
	helpers.QuitWitMessageIfErr(err, "Unable to parse JSON in the codes file. Quitting.")

	for _, code := range codes {
		fmt.Println(code)
	}

}
