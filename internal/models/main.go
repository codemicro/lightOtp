package models

import (
	"github.com/codemicro/lightOtp/internal/helpers"
	"path"
)

type TOTPCode struct {
	Issuer      string `json:"issuer"`
	AccountName string `json:"accountName"`
	Digits      int    `json:"digits"`
	Secret      string `json:"secret"`
}

type Settings struct {
	CodesLocation     string `json:"codesLocation"`
	DefaultCodeLength int    `json:"defaultCodeLength"`
}

func NewSettings() (Settings, error) {

	configDir, err := helpers.SetupConfigDir()
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		CodesLocation:     path.Join(configDir, "codes.json"),
		DefaultCodeLength: 6,
	}, nil

}
