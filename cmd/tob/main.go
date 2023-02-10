package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/notificators/discord"
)

func main() {
	args, err := tob.ParseArgument()
	if err != nil {
		fmt.Println("error: ", err)
		args.Help()
		os.Exit(1)
	}

	if args.ShowVersion {
		fmt.Printf("%s version %s\n", os.Args[0], tob.Version)
		os.Exit(0)
	}

	configFile, err := os.Open(args.ConfigFile)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	configs, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	// close configFile
	err = configFile.Close()
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	// init Notificator
	discordNotificator, err := discord.NewDiscord(configs)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	checkIntervalF, ok := configs["checkInterval"].(float64)
	if !ok {
		fmt.Println("error: cannot find checkInterval in config file")
		os.Exit(1)
	}

	// convert to int
	checkInterval := int(checkIntervalF)
	fmt.Println(checkInterval)

	if checkInterval <= 0 {
		// set check interval to 5 minutes
		checkInterval = 5000
	}

	notificators := []tob.Notificator{discordNotificator}

	// runner
	runner, err := tob.NewRunner(checkInterval, notificators, configs, args.Verbose)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	// initialize service from runner
	runner.InitServices()

	fmt.Println("wait...")

	kill := make(chan os.Signal, 1)
	// notify when user interrupt the process
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	go waitNotify(kill, runner)

	// run the Runner
	runner.Run()

}

func waitNotify(kill chan os.Signal, runner *tob.Runner) {
	select {
	case <-kill:
		runner.Stop() <- true
		fmt.Println("kill process")
	}
}
