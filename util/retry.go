package util

import (
	"context"
	"time"
)

// RetryFunc is a function type that will try to be called again if it doesn't work
type RetryFunc func() bool

// RetryCallbackFunc is a function type for callback
type RetryCallbackFunc func()

// Retry is a wrapper function that is used to try to call external resources that have the potential to fail
func Retry(ctx context.Context, n int, interval time.Duration,
	done chan<- bool, f RetryFunc, onSucceed RetryCallbackFunc, onError RetryCallbackFunc) {
	var i int = 1
	go func(ctx context.Context, i int, n int, interval time.Duration,
		done chan<- bool, f RetryFunc, onSucceed RetryCallbackFunc, onError RetryCallbackFunc) {

		for {
			select {
			case <-ctx.Done():
				done <- true
				return
			default:
			}

			succeed := f()
			if succeed {
				done <- true

				if onSucceed != nil {
					onSucceed()
				}
				return
			}

			if i == n {
				done <- true

				if onError != nil {
					onError()
				}
				return
			}

			<-time.After(interval)

			i++

		}
	}(ctx, i, n, interval, done, f, onSucceed, onError)
}

// Usage

// func main() {
// 	f := func() bool {
// 		resp, err := httpPost("not_ok")
// 		if err != nil {
// 			fmt.Println(err)
// 			return true
// 		}

// 		statusOK := resp.httpCode >= 200 && resp.httpCode < 300
// 		if !statusOK {
// 			fmt.Println("retry")
// 			return false
// 		}

// 		return true
// 	}

// 	done := make(chan bool, 1)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	//ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
// 	defer func() { cancel() }()

// 	onSucceed := func() {
// 		fmt.Println("execution succeed")
// 	}

// 	onError := func() {
// 		fmt.Println("execution error")
// 	}

// 	Retry(ctx, 5, 3*time.Second, done, f, onSucceed, onError)

// 	<-done

// 	//time.Sleep(time.Second * 15)
// }

// type response struct {
// 	httpCode int
// }

// func httpPost(data string) (*response, error) {
// 	if data == "fatal" {
// 		return nil, errors.New("fatal error")
// 	}

// 	if data == "not_ok" {
// 		return &response{httpCode: 400}, nil
// 	}

// 	return &response{httpCode: 200}, nil
// }
