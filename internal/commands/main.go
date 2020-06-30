package commands

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codemicro/lightOtp/internal/helpers"
	"github.com/codemicro/lightOtp/internal/persist"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

func Help() {
	fmt.Println("lightOTP help\n  help: shows this command\n  list: lists all providers added\n  code: gets " +
		"code for a provider (args: provider id)\n  new : adds a new provider\n  exit: I honestly have no clue... do" +
		" you?")
}

func ListProviders() {
	if len(persist.Codes) == 0 {
		helpers.PrintErrLn("No providers added")
		return
	}

	fmt.Println("Available providers:")

	for i, v := range persist.Codes {

		var accountName string
		if v.AccountName != "" {
			accountName = " (" + v.AccountName + ")"
		} else {
		}

		fmt.Printf("  %v: %s%s\n", i+1, v.Issuer, accountName)
	}
}

func GenerateCode(id int32) {
	code, err := totp.GenerateCodeCustom(persist.Codes[id].Secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.Digits(persist.Codes[id].Digits),
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		helpers.PrintErrLn("Unable to generate code - " + err.Error())
		return
	}

	var message string

	message += persist.Codes[id].Issuer
	if persist.Codes[id].AccountName != "" {
		message += " ("
		message += persist.Codes[id].AccountName
		message += ")"
	}
	message += ": "
	message += code

	fmt.Println(message)

	currentSecond := time.Now().Second()
	var remainingSeconds int
	if currentSecond < 30 {
		remainingSeconds = 30 - currentSecond
	} else {
		remainingSeconds = 60 - currentSecond
	}

	fmt.Printf("Valid for %d seconds\n", remainingSeconds)

	// Send code to clipboard

	if !clipboard.Unsupported {
		err = clipboard.WriteAll(code)
		if err != nil {
			helpers.PrintErrLn("Unable to copy code to clipboard")
		} else {
			fmt.Println("Copied code to clipboard")
		}
	}
}
