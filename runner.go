package tob

import (
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/dummy"
	"github.com/telkomdev/tob/postgres"
	"net/url"
	"time"
)

// Runner the tob runner
type Runner struct {
	configs         config.Config
	services        map[string]Service
	stopChan        chan bool
	verbose         bool
	intervalSeconds int
	initialized     bool
	notificators    []Notificator
}

// NewRunner Runner's constructor
func NewRunner(intervalSeconds int, notificators []Notificator, configs config.Config, verbose bool) (*Runner, error) {
	runner := new(Runner)
	runner.configs = configs

	stopChan := make(chan bool, 1)
	runner.stopChan = stopChan

	services := make(map[string]Service)
	runner.services = services

	runner.verbose = verbose

	runner.intervalSeconds = intervalSeconds
	runner.notificators = notificators

	return runner, nil
}

// Add will add new service to Runner
func (r *Runner) Add(service Service) {
	if service != nil {
		r.services[service.Name()] = service
	}
}

// InitServices will init initial services
func (r *Runner) InitServices() error {
	dummyService := dummy.NewDummy(r.verbose, Logger)
	r.Add(dummyService)

	postgresService := postgres.NewPostgres(r.verbose, Logger)
	r.Add(postgresService)

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

		fmt.Println(name, " ", conf["url"])

		if service, ok := r.services[name]; ok {
			urlStr, _ := conf["url"].(string)

			// validate and parse urlStr
			_, err := url.Parse(urlStr)
			if err != nil {
				return err
			}

			service.SetURL(urlStr)

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
func (r *Runner) Run() {
	if !r.initialized {
		panic("service not initialized yet")
	}

	if r.notificators == nil || len(r.notificators) <= 0 {
		panic("notificator cannot be nil")
	}

	ticker := time.NewTicker(time.Second * time.Duration(r.intervalSeconds))

	for {
		select {
		case <-r.stopChan:
			Logger.Println("runner received stop channel, cleanup resource now !!")
			r.cleanup()

			return
		case <-ticker.C:
			r.healthCheck()
			// Logger.Println(fmt.Sprintf("ticked: %s", t.String()))
		}
	}

}

// healthCheck will check health of all services
func (r *Runner) healthCheck() {
	for _, service := range r.services {
		go func(s Service) {
			resp := s.Ping()
			respStr := string(resp)
			if respStr == NotOk && s.IsRecover() {
				// set recover to false
				s.SetRecover(false)

				for _, notificator := range r.notificators {
					err := notificator.Send(fmt.Sprintf("%s is DOWN", s.Name()))
					if err != nil {
						Logger.Println(err)
					}
				}
			}

			if respStr == OK && !s.IsRecover() {
				// set recover to true
				s.SetRecover(true)

				for _, notificator := range r.notificators {
					err := notificator.Send(fmt.Sprintf("%s is UP", s.Name()))
					if err != nil {
						Logger.Println(err)
					}
				}
			}

			Logger.Println(fmt.Sprintf("%s => %s", s.Name(), respStr))
		}(service)
	}
}

// Stop will receive stop channel
func (r *Runner) Stop() chan<- bool {
	return r.stopChan
}

// cleanup will Cleanup the tob Runner services resource
func (r *Runner) cleanup() error {
	for _, service := range r.services {
		err := service.Close()
		if err != nil {
			Logger.Println(err)
		}
	}

	return nil
}
