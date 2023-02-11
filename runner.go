package tob

import (
	"context"
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/services/dummy"
	"github.com/telkomdev/tob/services/postgres"
	"github.com/telkomdev/tob/services/web"
	"net/url"
	"time"
)

// Runner the tob runner
type Runner struct {
	configs      config.Config
	services     map[string]Service
	stopChan     chan bool
	verbose      bool
	initialized  bool
	notificators []Notificator
}

// NewRunner Runner's constructor
func NewRunner(notificators []Notificator, configs config.Config, verbose bool) (*Runner, error) {
	runner := new(Runner)
	runner.configs = configs

	stopChan := make(chan bool, 1)
	runner.stopChan = stopChan

	services := make(map[string]Service)
	runner.services = services

	runner.verbose = verbose

	runner.notificators = notificators

	return runner, nil
}

func initServiceKind(serviceKind ServiceKind, verbose bool) (Service, bool) {
	services := make(map[ServiceKind]Service)
	services[Dummy] = dummy.NewDummy(verbose, Logger)
	services[Postgresql] = postgres.NewPostgres(verbose, Logger)
	services[Web] = web.NewWeb(verbose, Logger)

	s, ok := services[serviceKind]
	return s, ok
}

// Add will add new service to Runner
func (r *Runner) Add(service Service) {
	if service != nil {
		r.services[service.Name()] = service
	}
}

// InitServices will init initial services
func (r *Runner) InitServices() error {
	serviceConfigInterface, ok := r.configs["service"]
	if !ok {
		return errors.New("field service not found in config file")
	}

	serviceConfigs, ok := serviceConfigInterface.(map[string]interface{})
	if !ok {
		return errors.New("invalid config file")
	}

	for name, confInterface := range serviceConfigs {
		conf, ok := confInterface.(map[string]interface{})
		if !ok {
			return errors.New("invalid config file")
		}

		Logger.Println(name, " ", conf["url"])

		urlStr, ok := conf["url"].(string)
		if !ok {
			return errors.New("invalid config file")
		}

		serviceKind, ok := conf["kind"].(string)
		if !ok {
			return errors.New("invalid config file")
		}

		checkIntervalF, ok := conf["checkInterval"].(float64)
		if !ok {
			return errors.New("invalid config file")
		}

		// convert to int
		checkInterval := int(checkIntervalF)

		// set default checkInterval
		if checkInterval <= 0 {
			// set check interval to 5 minutes
			checkInterval = 5000
		}

		serviceEnabled, ok := conf["enable"].(bool)
		if !ok {
			return errors.New("invalid config file")
		}

		if s, ok := initServiceKind(ServiceKind(serviceKind), r.verbose); ok {
			r.services[name] = s
		}

		if service, ok := r.services[name]; ok && service != nil && serviceEnabled {

			// validate and parse urlStr
			_, err := url.Parse(urlStr)
			if err != nil {
				return err
			}

			service.SetURL(urlStr)
			service.SetCheckInterval(checkInterval)
			service.Enable(serviceEnabled)

			err = service.Connect()
			if err != nil {
				return err
			}
		}
	}

	// set initialized to true
	r.initialized = true

	return nil
}

// Run will Run the tob Runner
func (r *Runner) Run(ctx context.Context) {
	if !r.initialized {
		panic("service not initialized yet")
	}

	if r.notificators == nil || len(r.notificators) <= 0 {
		panic("notificator cannot be nil")
	}

	totalServiceToBeExecuted := 0

	for name, service := range r.services {
		if service != nil && service.IsEnabled() {
			totalServiceToBeExecuted++

			ticker := time.NewTicker(time.Second * time.Duration(service.GetCheckInterval()))

			go func(n string, s Service, t *time.Ticker) {
				for {
					select {
					case <-s.Stop():
						Logger.Println(fmt.Sprintf("runner service %s received stop channel, cleanup resource now !!", n))

						return
					case <-t.C:
						resp := s.Ping()
						respStr := string(resp)
						if respStr == NotOk && s.IsRecover() {
							// set recover to false
							s.SetRecover(false)

							for _, notificator := range r.notificators {
								err := notificator.Send(fmt.Sprintf("%s is DOWN", n))
								if err != nil {
									Logger.Println(err)
								}
							}
						}

						if respStr == OK && !s.IsRecover() {
							// set recover to true
							s.SetRecover(true)

							for _, notificator := range r.notificators {
								err := notificator.Send(fmt.Sprintf("%s is UP", n))
								if err != nil {
									Logger.Println(err)
								}
							}
						}

						Logger.Println(fmt.Sprintf("%s => %s", n, respStr))
					}
				}
			}(name, service, ticker)

		}
	}

	if r.verbose {
		Logger.Println(fmt.Sprintf("total service to be executed: %d", totalServiceToBeExecuted))
	}

	// block here
	for {
		// The try-receive operation here is to
		// try to exit the worker goroutine as
		// early as possible. Try-receive
		// optimized by the standard Go
		// compiler, so they are very efficient.
		select {
		case <-ctx.Done():
			Logger.Println("runner context canceled")
			r.cleanup()
			return
		default:
		}

		select {
		case <-r.stopChan:
			Logger.Println("runner received stop channel, cleanup resource now !!")
			r.cleanup()
			return

		case <-ctx.Done():
			Logger.Println("runner context canceled")
			r.cleanup()
			return
		}
	}

}

// Stop will receive stop channel
func (r *Runner) Stop() chan<- bool {
	return r.stopChan
}

// cleanup will Cleanup the tob Runner services resource
func (r *Runner) cleanup() error {
	for _, service := range r.services {
		if service != nil && service.IsEnabled() {
			err := service.Close()
			if err != nil {
				Logger.Println(err)
			}

			// send stop channel
			service.Stop() <- true
		}
	}

	return nil
}
