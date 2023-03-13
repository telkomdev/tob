package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/telkomdev/tob"
)

const (
	// DefaultPort default port for tob-http-agent
	DefaultPort = 9113

	// Version current version for tob-http-agent
	Version = "1.0.0"
)

func main() {

	var (
		httpPort    int
		showVersion bool
		verbose     bool
	)

	flag.IntVar(&httpPort, "port", DefaultPort, "HTTP Port (eg: 9113)")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&verbose, "V", true, "verbose mode (if true log will appear otherwise no)")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: ")
		fmt.Println("tob-http-agent -[options]")
		fmt.Println()
		fmt.Println("tob-http-agent -port 9113")
		fmt.Println()
		fmt.Println("-h | -help (show help)")
		fmt.Println("-v | -version (show version)")
		fmt.Println("-V : verbose mode")
		fmt.Println("---------------------------")
		fmt.Println()
	}

	flag.Parse()

	// show version
	if showVersion {
		fmt.Printf("%s version %s (runtime: %s)\n", os.Args[0], Version, runtime.Version())
		os.Exit(0)
	}

	http.HandleFunc("/", loggerMiddleware(indexHandler()))
	http.HandleFunc("/check-disk", loggerMiddleware(checkStorageHandler()))

	tob.Logger.Printf("webapp running on port :%d\n", httpPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
	if err != nil {
		tob.Logger.Println(err)
		os.Exit(1)
	}
}

type customResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func indexHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		jsonResponse(res, 200, []byte(`{"success": true, "message": "server up and running"}`))
	}
}

func checkStorageHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		jsonMap, wait, err := checkDiskStatus()
		if err != nil {
			jsonResponse(res, 500, []byte(`{"success": false, "message": "error check storage"}`))
			return
		}

		err = wait()
		if err != nil {
			tob.Logger.Println(err)
			jsonResponse(res, 500, []byte(`{"success": false, "message": "error check storage"}`))
			return
		}

		payload := customResponse{
			Success: true,
			Message: "disk status",
			Data:    jsonMap,
		}

		response, err := json.Marshal(payload)

		if err != nil {
			tob.Logger.Println(err)
			jsonResponse(res, 500, []byte(`{"success": false, "message": "error check storage"}`))
			return
		}

		jsonResponse(res, 200, response)
	}
}

// HTTP utility
func jsonResponse(res http.ResponseWriter, httpCode int, payload []byte) {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(payload)
}

func loggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		tob.Logger.Printf("path: %s | method: %s | remote_address: %s | user_agent: %s | duration: %v",
			req.URL.EscapedPath(), req.Method, req.RemoteAddr, req.UserAgent(), time.Since(start))
		next.ServeHTTP(res, req)
	}
}

func checkDiskStatus() (map[string]interface{}, func() error, error) {
	dfPath, err := exec.LookPath("df")
	if err != nil {
		tob.Logger.Println(err)
		return nil, nil, err
	}

	cmd := exec.Command(dfPath, "/")
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		tob.Logger.Println(err)
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		tob.Logger.Println(err)
		return nil, nil, err
	}

	b, err := io.ReadAll(outPipe)
	if err != nil {
		tob.Logger.Println(err)
		return nil, nil, err
	}

	ss := strings.Split(string(b), "\n")

	var (
		jsonMap map[string]interface{}
	)

	split := func(text string) ([]string, error) {
		re, err := regexp.Compile(`\s+|\s+$`)
		if err != nil {
			return nil, err
		}
		result := re.Split(text, -1)
		return result, nil
	}

	parseNumberOnly := func(str string) int {
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		allStringNum := re.FindAllString(str, -1)
		if len(allStringNum) > 0 {
			numStr := allStringNum[0]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return 0
			}

			return num
		}
		return 0
	}

	if len(ss) >= 2 {
		headers := ss[0]
		values := ss[1]

		tob.Logger.Println(headers)
		tob.Logger.Println(values)

		headerSplitted, err := split(headers)
		if err != nil {
			tob.Logger.Println(err)
			return nil, nil, err
		}

		valueSplitted, err := split(values)
		if err != nil {
			tob.Logger.Println(err)
			return nil, nil, err
		}

		// the length of headerSplitted must be equal to valueSplitted
		headerSplitted = headerSplitted[:len(valueSplitted)]

		jsonMap = make(map[string]interface{})

		for i := 0; i < len(headerSplitted); i++ {
			if strings.ToLower(headerSplitted[i]) == "use%" {
				jsonMap[strings.ToLower(headerSplitted[i])] = parseNumberOnly(valueSplitted[i])
			} else {
				jsonMap[strings.ToLower(headerSplitted[i])] = valueSplitted[i]
			}

		}

	}

	wait := func() error {
		return cmd.Wait()
	}

	return jsonMap, wait, nil
}
