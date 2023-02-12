package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
)

func main() {
	args, err := tob.ParseArgument()
	if err != nil {
		fmt.Println("error: ", err)
		args.Help()
		os.Exit(1)
	}

	if args.ShowVersion {
		fmt.Printf("%s version %s (runtime: %s)\n", os.Args[0], tob.Version, runtime.Version())
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
	notificators, err := tob.InitNotificatorFactory(configs)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() { cancel() }()

	// runner
	runner, err := tob.NewRunner(notificators, configs, args.Verbose)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	// initialize service from runner
	err = runner.InitServices()
	if err != nil {
		fmt.Println("error InitServices : ", err)
		os.Exit(1)
	}

	kill := make(chan os.Signal, 1)
	// notify when user interrupt the process
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	go waitNotify(kill, runner)

	// run the Runner
	runner.Run(ctx)

}

func waitNotify(kill chan os.Signal, runner *tob.Runner) {
	select {
	case <-kill:
		runner.Stop() <- true
		fmt.Println("kill tob")
	}
}
