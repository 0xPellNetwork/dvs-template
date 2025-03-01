package tools

import (
	"flag"
	"fmt"
	"os"
)

// ParseCommandLineArgs parses command line arguments and returns command and home directory
// Returns command string and home directory path
func ParseCommandLineArgs() (string, *string) {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go [start-operator|start-taskgateway] --home=<home_dir>")
		os.Exit(1)
	}

	command := args[0]
	if command != "start-operator" && command != "start-task-gateway" {
		fmt.Println("Invalid command. Must be either 'start-app' or 'start-taskgateway'")
		os.Exit(1)
	}

	// Remove first argument (command) before parsing flags
	os.Args = append([]string{os.Args[0]}, args[1:]...)
	homeFlag := flag.String("home", "", "Node home directory")
	flag.Parse()

	if *homeFlag == "" {
		fmt.Println("--home flag is required")
		os.Exit(1)
	}

	return command, homeFlag
}
