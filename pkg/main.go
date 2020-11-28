package main

import (
	"flag"
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/config"
	"os"
	"strings"

	"github.com/AntonBorzenko/RestrictedPassportService/services"
)

func printHelp() {
	fmt.Println("api or update command is required")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	config.LoadConfigFromArgs(os.Args[2:])

	passportService := services.NewPassportService()

	switch strings.ToLower(os.Args[1]) {
	case "api":
		passportService.StartApi()
	case "update":
		passportService.UpdateDb()
	default:
		printHelp()
		flag.PrintDefaults()
		os.Exit(1)
	}
}
