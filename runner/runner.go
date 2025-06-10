package runner

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"plugin"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/services/airflow"
	"github.com/telkomdev/tob/services/diskstatus"
	"github.com/telkomdev/tob/services/dummy"
	"github.com/telkomdev/tob/services/elasticsearch"
	"github.com/telkomdev/tob/services/kafka"
	"github.com/telkomdev/tob/services/mongodb"
	"github.com/telkomdev/tob/services/mysqldb"
	"github.com/telkomdev/tob/services/oracle"
	"github.com/telkomdev/tob/services/postgres"
	"github.com/telkomdev/tob/services/redisdb"
	"github.com/telkomdev/tob/services/sslstatus"
	"github.com/telkomdev/tob/services/web"
	"github.com/telkomdev/tob/util"
)

// Runner the tob runner
type Runner struct {
	configs     config.Config
	services    map[string]tob.Service
	stopChan    chan bool
	verbose     bool
	initialized bool
	waiter      tob.Waiter
}

// NewRunner Runner's constructor
func NewRunner(configs config.Config, verbose bool) (*Runner, error) {
	runner := new(Runner)
	runner.configs = configs

	stopChan := make(chan bool, 1)
	runner.stopChan = stopChan

	services := make(map[string]tob.Service)
	runner.services = services

	runner.verbose = verbose

	return runner, nil
}

func lookupPlugin(pluginPath string) (tob.Service, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	serviceSymbol, err := plug.Lookup("Service")
	if err != nil {
		return nil, err
	}

	s, ok := serviceSymbol.(tob.Service)
	if !ok {
		return nil, errors.New("symbol is not valid tob.Service")
	}

	return s, nil
}

func initServiceKind(serviceKind tob.ServiceKind, pluginPath string, verbose bool) (tob.Service, bool) {
	services := make(map[tob.ServiceKind]tob.Service)
	services[tob.Airflow] = airflow.NewAirflow(verbose, tob.Logger)
	services[tob.AirflowFlower] = airflow.NewAirflowFlower(verbose, tob.Logger)
	services[tob.Dummy] = dummy.NewDummy(verbose, tob.Logger)
	services[tob.DiskStatus] = diskstatus.NewDiskStatus(verbose, tob.Logger)
	services[tob.Kafka] = kafka.NewKafka(verbose, tob.Logger)
	services[tob.MongoDB] = mongodb.NewMongo(verbose, tob.Logger)
	services[tob.MySQL] = mysqldb.NewMySQL(verbose, tob.Logger)
	services[tob.Postgresql] = postgres.NewPostgres(verbose, tob.Logger)
	services[tob.Oracle] = oracle.NewOracle(verbose, tob.Logger)
	services[tob.Redis] = redisdb.NewRedis(verbose, tob.Logger)
	services[tob.Web] = web.NewWeb(verbose, tob.Logger)
	services[tob.SSLStatus] = sslstatus.NewSSLStatus(verbose, tob.Logger)
	services[tob.Elasticsearch] = elasticsearch.NewElasticsearch(verbose, tob.Logger)

	if pluginPath != "" {
		servicePlugin, err := lookupPlugin(pluginPath)
		if err != nil {
			panic(err)
		}

		services[tob.Plugin] = servicePlugin
	}

	s, ok := services[serviceKind]
	return s, ok
}

// Add will add new service to Runner
func (r *Runner) Add(service tob.Service) {
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

	totalServiceToBeExecuted := 0

	for name, confInterface := range serviceConfigs {
		conf, ok := confInterface.(map[string]interface{})
		if !ok {
			return errors.New("invalid config file")
		}

		tob.Logger.Println(name)

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

		pluginPath, ok := conf["pluginPath"].(string)
		if !ok {
			pluginPath = ""
		}

		serviceEnabled, ok := conf["enable"].(bool)
		if !ok {
			return errors.New("invalid config file")
		}

		if serviceEnabled {
			if s, ok := initServiceKind(tob.ServiceKind(serviceKind), pluginPath, r.verbose); ok {
				r.services[name] = s
			}
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
			service.SetConfig(conf)
			service.SetNotificatorConfig(conf)

			// by default service is recovered
			service.SetRecover(true)

			err = service.Connect()
			if err != nil {
				return err
			}

			totalServiceToBeExecuted++
		}
	}

	// set initialized to true
	r.initialized = true

	// set waiter capacity with amount of service to be executed
	r.waiter = tob.NewWaiter(uint(totalServiceToBeExecuted))
	if r.verbose {
		tob.Logger.Printf("total service to be executed: %d\n", uint(totalServiceToBeExecuted))
	}

	return nil
}

func healthCheck(n string, s tob.Service, t *time.Ticker, waiter tob.Waiter) {

	for {
		select {
		case <-s.Stop():
			tob.Logger.Printf("runner service %s received stop channel, cleanup resource now !!\n", n)

			// stop ticker
			t.Stop()

			// tell waiter this service execution is done
			waiter.Done()

			return
		case <-t.C:
			// set message to empty
			s.SetMessage("")

			resp := s.Ping()
			respStr := string(resp)

			// Airflow Monitoring
			if s.Name() == string(tob.Airflow) {
				for _, notificator := range s.GetNotificators() {
					if !util.IsNilish(notificator) {
						if notificator.IsEnabled() && notificator.Provider() == "webhook" {
							notificatorMessage := fmt.Sprintf("%s is DOWN", n)
							if s.GetMessage() != "" {
								notificatorMessage = fmt.Sprintf("%s is CHECKING | %s", n, s.GetMessage())
								if respStr == tob.NotOk {
									notificatorMessage = fmt.Sprintf("%s is DOWN | %s", n, s.GetMessage())
								}
							}
							if notificator.IsEnabled() && s.Name() != string(tob.SSLStatus) {
								err := notificator.Send(notificatorMessage)
								if err != nil {
									tob.Logger.Printf("notificator %s error: %s", notificator.Provider(), err.Error())
								}
							}
						}
					}
				}
			}

			// SSL Monitoring
			if s.Name() == string(tob.SSLStatus) {
				for _, notificator := range s.GetNotificators() {
					if !util.IsNilish(notificator) {
						if notificator.IsEnabled() && notificator.Provider() == "webhook" {
							notificatorMessage := fmt.Sprintf("%s is DOWN", n)
							if s.GetMessage() != "" {
								notificatorMessage = fmt.Sprintf("%s is MONITORED | %s", n, s.GetMessage())
							}
							err := notificator.Send(notificatorMessage)
							if err != nil {
								tob.Logger.Printf("notificator %s error: %s", notificator.Provider(), err.Error())
							}
						}
					}
				}
			}

			if respStr == tob.NotOk && s.IsRecover() {
				// set last downtime
				s.SetLastDownTimeNow()
				// set recover to false
				s.SetRecover(false)

				notificatorMessage := fmt.Sprintf("%s is DOWN", n)
				if s.GetMessage() != "" {
					notificatorMessage = fmt.Sprintf("%s is DOWN | %s", n, s.GetMessage())
				}

				for _, notificator := range s.GetNotificators() {
					if !util.IsNilish(notificator) {
						if notificator.IsEnabled() && s.Name() != string(tob.SSLStatus) {
							err := notificator.Send(notificatorMessage)
							if err != nil {
								tob.Logger.Printf("notificator %s error: %s", notificator.Provider(), err.Error())
							}
						}
					}
				}

			}

			if respStr == tob.OK && !s.IsRecover() {
				// set recover to true
				s.SetRecover(true)

				notificatorMessage := fmt.Sprintf("%s is UP. It was down for %s", n, s.GetDownTimeDiff())
				if s.GetMessage() != "" {
					notificatorMessage = fmt.Sprintf("%s is UP | %s", n, s.GetMessage())
				}

				for _, notificator := range s.GetNotificators() {
					if !util.IsNilish(notificator) {
						if notificator.IsEnabled() && s.Name() != string(tob.SSLStatus) {
							err := notificator.Send(notificatorMessage)
							if err != nil {
								tob.Logger.Printf("notificator %s error: %s\n", notificator.Provider(), err.Error())
							}
						}
					}
				}
			}

			tob.Logger.Printf("%s => %s\n", n, respStr)
		}
	}
}

// Run will Run the tob Runner
func (r *Runner) Run(ctx context.Context) {
	if !r.initialized {
		panic("service not initialized yet")
	}

	// close waiter's channel indicates that no more values will be sent on it
	defer func() { r.waiter.Close() }()

	for name, service := range r.services {
		if service != nil && service.IsEnabled() {

			ticker := time.NewTicker(time.Second * time.Duration(service.GetCheckInterval()))

			// run all services health check on its goroutine
			go healthCheck(name, service, ticker, r.waiter)

		}
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
			tob.Logger.Println("runner context canceled")
			r.cleanup()

			// wait all service's goroutine to stop
			r.waiter.Wait()
			return
		default:
		}

		select {
		case <-r.stopChan:
			tob.Logger.Println("runner received stop channel, cleanup resource now !!")
			r.cleanup()

			// wait all service's goroutine to stop
			r.waiter.Wait()
			return

		case <-ctx.Done():
			tob.Logger.Println("runner context canceled")
			r.cleanup()

			// wait all service's goroutine to stop
			r.waiter.Wait()
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
				tob.Logger.Println(err)
			}

			// send stop channel
			service.Stop() <- true
		}
	}

	return nil
}
