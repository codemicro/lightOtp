package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codemicro/lightOtp/internal/helpers"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/codemicro/lightOtp/internal/persist"
	"github.com/fatih/color"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"io/ioutil"
	"net/url"
	"os"
	"time"
)

func Help() {
	fmt.Println("lightOTP help\n  help: shows this command\n  list: lists all providers added\n  code: gets " +
		"code for a provider (args: provider id)\n  add : adds a new provider\n  exit: ¯\\_(ツ)_/¯")
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

func AddProvider() {
	scanner := bufio.NewScanner(os.Stdin)
	_, _ = color.New(color.FgCyan).Print("Secret or URI > ")
	scanner.Scan()

	text := scanner.Text()

	var issuer string
	var accountName string
	var secret string
	digits := 6

	_, err := url.ParseRequestURI(text)
	if err == nil {
		// URI string is valid
		fromUri, err := otp.NewKeyFromURL(text)
		if err != nil {
			helpers.ErrWithMessage(err, "Unable to parse URI")
			return
		}

		issuer = fromUri.Issuer()
		accountName = fromUri.AccountName()
		secret = fromUri.Secret()
	} else {
		// Must be a plaintext secret - prompt for account name and issuer
		secret = text

		for issuer == "" {
			_, _ = color.New(color.FgCyan).Print("Provider name > ")
			scanner.Scan()
			issuer = scanner.Text()
		}

		_, _ = color.New(color.FgCyan).Print("Account name (default '') > ")
		scanner.Scan()
		accountName = scanner.Text()

	}

	codeInst := models.TOTPCode{
		Issuer:      issuer,
		AccountName: accountName,
		Digits:      digits,
		Secret:      secret,
	}

	persist.Codes = append(persist.Codes, codeInst)

	jsonCodes, _ := json.Marshal(persist.Codes)

	err = ioutil.WriteFile(persist.Settings.CodesLocation, jsonCodes, 0644)
	if err != nil {
		helpers.PrintErrLn("Unable to save codes to file.")
	} else {
		fmt.Printf("Added as ID %d", len(persist.Codes))
	}

}
