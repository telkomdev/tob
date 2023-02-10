package tob

import (
	"flag"
	"fmt"
)

// Argument type
type Argument struct {
	ShowVersion bool
	ConfigFile  string
	Help        func()
	Message     []byte
	Verbose     bool
}

// ParseArgument will parse the OS args into Argument type
func ParseArgument() (*Argument, error) {

	argument := new(Argument)

	var (
		configFile  string
		showVersion bool
		verbose     bool
	)

	flag.StringVar(&configFile, "config", "config.json", "config file (.json)")
	flag.StringVar(&configFile, "c", "config.json", "config file (.json)")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&verbose, "V", true, "verbose mode (if true log will appear otherwise no)")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: ")
		fmt.Println("tob -[options]")
		fmt.Println()
		fmt.Println("-config | -c (configuration .json file)")
		fmt.Println("-h | -help (show help)")
		fmt.Println("-v | -version (show version)")
		fmt.Println("-V : verbose mode")
		fmt.Println("---------------------------")
		fmt.Println()
	}

	flag.Parse()

	argument.ConfigFile = configFile
	argument.ShowVersion = showVersion
	argument.Verbose = verbose
	return argument, nil
}
