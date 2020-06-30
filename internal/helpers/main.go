package helpers

import (
	"encoding/json"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/codemicro/lightOtp/internal/persist"
	"io/ioutil"
	"os"
	"path"
)

func OpenConfigFile(fname string) (string, string, error) {

	// Returns as "" if empty or missing

	directoryPath, err := SetupConfigDir()

	if err != nil {
		return "", "", err
	}

	fLoc := path.Join(directoryPath, fname)

	fileConts, err := ioutil.ReadFile(fLoc)

	if err != nil {
		return "", "", nil
	}

	return string(fileConts), fLoc, nil

}

func NewSettings() (models.Settings, error) {

	configDir, err := SetupConfigDir()
	if err != nil {
		return models.Settings{}, err
	}

	return models.Settings{
		CodesLocation:     path.Join(configDir, "codes.json"),
		DefaultCodeLength: 6,
	}, nil

}

func SetupConfigDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	directoryPath := path.Join(userHomeDir, ".config", "lightOtp")

	_ = os.Mkdir(directoryPath, os.ModeDir) // Ignore error (thrown when dir already exists)

	return directoryPath, nil
}

func UpdateCodes() error {
	jsonCodes, _ := json.Marshal(persist.Codes)
	err := ioutil.WriteFile(persist.Settings.CodesLocation, jsonCodes, 0644)
	return err
}

func CheckErr(err error) {
	if err != nil {
		PrintErrLn(err.Error())
		os.Exit(1)
	}
}

func ErrWithMessage(err error, message string) {
	if err != nil {
		PrintErrLn(message)
		if os.Getenv("DEBUG") != "" {
			PrintDebugLn(err.Error())
		}
	}
}

func QuitWithMessageIfErr(err error, message string) {
	if err != nil {
		PrintErrLn(message)
		if os.Getenv("DEBUG") != "" {
			PrintDebugLn(err.Error())
		}
		os.Exit(1)
	}
}
