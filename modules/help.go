package modules

import (
	"awesomeProject/utils"
	"sync"
	"time"
)

var IsProxylessMode bool = false
var ProxyType string = "HTTP"
var Retries int = 0
var DebounceKey string = "sex"
var countryCode string = "33"
var DebounceIdentifier = "[" + time.Now().Format("01-02-2006") + " at " + time.Now().Format("15-04-05") + "].txt"
var Speed int64

type SafeChecker struct {
	mu sync.Mutex
	v  []string
}

func EmailCheck(emailSort func(email string, wg *sync.WaitGroup, proxy string, mu *sync.Mutex), emails []string, wg *sync.WaitGroup, proxies []string, mu *sync.Mutex)  {
	SC := SafeChecker{}
	SC.v = emails

	proxyOffset := 0
	proxyOffsetMax := len(proxies) - 1


	dvd := utils.Chunks(emails, int(3000))

	for i := 0; i < len(dvd); i++ {
		for _, email := range dvd[i] {
			wg.Add(1)
			go emailSort(email, wg, proxies[proxyOffset], mu)

			if !IsProxylessMode {
				if proxyOffset == proxyOffsetMax  {
					proxyOffset = 0
				} else {
					proxyOffset++
				}
			}

			if Speed == 50 {
				time.Sleep(time.Millisecond * 50)
			} else if Speed == 10 {
				time.Sleep(time.Millisecond * 20)
			} else if Speed == 1 {
				time.Sleep(time.Millisecond * 10)
			}



		}
		time.Sleep(time.Millisecond * 5000)
	}

	wg.Wait()
}