package commands

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/codemicro/lightOtp/internal/helpers"
	"github.com/codemicro/lightOtp/internal/models"
	"github.com/codemicro/lightOtp/internal/persist"
	"github.com/fatih/color"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"net/url"
	"os"
	"strings"
	"time"
)

func Help() {
	fmt.Println("lightOTP help\n  help: shows this command\n  list: lists all providers added\n  code: gets " +
		"code for a provider (args: provider ID)\n  add : adds a new provider\n  del : removes a provider (args: " +
		"provider ID)\n  exit: ¯\\_(ツ)_/¯")
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
		Period:    persist.Codes[id].Period,
		Skew:      1,
		Digits:    otp.Digits(persist.Codes[id].Digits),
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		helpers.ErrWithMessage(err, "Unable to generate code")
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
			helpers.ErrWithMessage(err, "Unable to copy code to clipboard")
		} else {
			fmt.Println("Copied code to clipboard")
		}
	}
}

func AddProvider() {
	// TODO: Validate secret is valid base32 thing

	scanner := bufio.NewScanner(os.Stdin)
	_, _ = color.New(color.FgCyan).Print("Secret or URI > ")

	text := helpers.CollectCensoredInput()

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
		Period:      30,
	}

	persist.Codes = append(persist.Codes, codeInst)
	err = helpers.UpdateCodes()
	if err != nil {
		helpers.ErrWithMessage(err, "Unable to save codes to file.")
	} else {
		fmt.Printf("Added as ID %d\n", len(persist.Codes))
	}
}

func RemoveProvider(id int32) {

	scanner := bufio.NewScanner(os.Stdin)
	_, _ = color.New(color.FgRed).Print("Are you sure? (y/N) > ")
	scanner.Scan()

	if strings.ToLower(scanner.Text()) == "y" {
		persist.Codes = append(persist.Codes[:id], persist.Codes[id+1:]...)

		err := helpers.UpdateCodes()
		if err != nil {
			helpers.ErrWithMessage(err, "Unable to save codes file.")
		} else {
			fmt.Println("Deleted.")
		}
	} else {
		fmt.Println("No action taken.")
	}

}
