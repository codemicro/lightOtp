package helpers

import (
	"fmt"
	"github.com/fatih/color"
)

// Printing help

func PrintErrLn(message string) {
	fmt.Print("[ ")
	_, _ = color.New(color.FgRed).Print("ERROR")
	fmt.Println(" ] " + message)
}

func PrintInfoLn(message string) {
	fmt.Print("[ ")
	_, _ = color.New(color.FgYellow).Print("INFO ")
	fmt.Println(" ] " + message)
}

func PrintDebugLn(message string) {
	fmt.Print("[ ")
	_, _ = color.New(color.FgRed).Print("D")
	_, _ = color.New(color.FgYellow).Print("E")
	_, _ = color.New(color.FgGreen).Print("B")
	_, _ = color.New(color.FgBlue).Print("U")
	_, _ = color.New(color.FgMagenta).Print("G")
	fmt.Println(" ] " + message)
}