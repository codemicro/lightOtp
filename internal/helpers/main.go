package helpers

import (
	"encoding/json"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/codemicro/lightOtp/internal/persist"
	"io/ioutil"
	"os"
	"path"
)

var(
	ConfigDirectory = func() string {
		userHomeDir, err := os.UserHomeDir()
		QuitWithMessageIfErr(err, "Unable to locate user's home directory.")
		return path.Join(userHomeDir, ".config", "lightOtp")
	}()
)

func OpenConfigFile(fname string) (string, string, error) {

	// Returns as "" if empty or missing

	SetupConfigDir()
	fLoc := path.Join(ConfigDirectory, fname)

	fileConts, err := ioutil.ReadFile(fLoc)

	if err != nil {
		return "", "", nil
	}

	return string(fileConts), fLoc, nil

}

func NewSettings() (models.Settings, error) {

	SetupConfigDir()

	return models.Settings{
		CodesLocation:     path.Join(ConfigDirectory, "codes.json"),
		DefaultCodeLength: 6,
	}, nil

}

func SetupConfigDir() {
	_ = os.Mkdir(ConfigDirectory, os.ModeDir) // Ignore error (thrown when dir already exists)
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
