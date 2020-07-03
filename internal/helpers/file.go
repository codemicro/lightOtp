package helpers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/codemicro/lightOtp/internal/encryption"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/codemicro/lightOtp/internal/persist"
	"io/ioutil"
)

func LoadCodes() error {
	// Parse wrapper JSON
	wrapperJsonSlice, err := ioutil.ReadFile(persist.Settings.CodesLocation)
	var wrapperJson models.CodesFile
	err = json.Unmarshal(wrapperJsonSlice, &wrapperJson)
	if err != nil {
		return err
	}

	// Decode and decrypt codes JSON
	encryptedCodes, err := base64.StdEncoding.DecodeString(wrapperJson.Codes)
	codesJsonSlice, err := encryption.Decrypt(encryptedCodes, persist.MasterPassword)
	if err != nil {
		return err
	}
	// Validate checksum
	checksum := md5.Sum(codesJsonSlice)
	checksumString := encryption.ConvertBytesToHex(checksum[:])
	if wrapperJson.Checksum != checksumString {
		return errors.New("checksum does not match")
	}
	// Store codes JSON into variable
	err = json.Unmarshal(codesJsonSlice, &persist.Codes)
	return err
}

func UpdateCodes() error {
	// Generate codes JSON text
	jsonCodes, _ := json.Marshal(persist.Codes)
	// Generate MD5 checksum
	checksum := md5.Sum(jsonCodes)
	checksumString := encryption.ConvertBytesToHex(checksum[:])
	// Encrypt JSON text
	cipherText, err := encryption.Encrypt(jsonCodes, persist.MasterPassword)
	if err != nil {
		return err
	}
	encodedCodes := base64.StdEncoding.EncodeToString(cipherText)
	// Save into JSON to write to disk
	codesStruct := models.CodesFile{
		Checksum: checksumString,
		Codes:    encodedCodes,
	}
	jsonToWrite, _ := json.Marshal(codesStruct)
	err = ioutil.WriteFile(persist.Settings.CodesLocation, jsonToWrite, 0644)
	return err
}
