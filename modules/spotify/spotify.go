package spotify

import (
	"awesomeProject/modules"
	"awesomeProject/utils"
	"github.com/pterm/pterm"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var good int64 = 0
var bad int64 = 0
var total int64 = 0
var error int64 = 0

func StartEmailChecker(emails []string, proxies []string, mode string)  {
	goodLog := pterm.NewStyle(pterm.FgGreen)
	badLog := pterm.NewStyle(pterm.FgRed)
	errLog := pterm.NewStyle(pterm.FgLightYellow)

	var debounceIdentifier = "Spotify" + modules.DebounceIdentifier

	_, err := os.Stat("./results/good/" + debounceIdentifier)
	_, err = os.Stat("./results/bad/" + debounceIdentifier)
	var goodFile *os.File
	var badFile *os.File

	if os.IsNotExist(err) {
		goodFile, _ = os.Create("./results/good/" + debounceIdentifier)
		badFile, _ = os.Create("./results/bad/" + debounceIdentifier)
	} else {
		goodFile, _ = os.Open("./results/good/" + debounceIdentifier)
		badFile, _ = os.Open("./results/bad/" + debounceIdentifier)
	}

	var goodMails []string
	var badMails []string
	var errorsMails []string
	var wg sync.WaitGroup
	var mu sync.Mutex



	modules.EmailCheck(
		func(email string, wg *sync.WaitGroup, proxyS string, mu *sync.Mutex) {
			defer wg.Done()
			httpTransport := &http.Transport{}

			if mode == "recheck" {
				proxyS = proxies[rand.Int() % len(proxies)]
			}

			if !modules.IsProxylessMode {
				if modules.ProxyType == "HTTP" {
					proxyUrl, _ := url.Parse("http://" + proxyS)
					httpTransport.Proxy = http.ProxyURL(proxyUrl)

				} else if modules.ProxyType == "SOCKS4" {
					dial := socks.Dial("socks4://" + proxyS)
					httpTransport.Dial = dial

				} else if modules.ProxyType == "SOCKS5" {
					dialer, err := proxy.SOCKS5("tcp", proxyS, nil, proxy.Direct)
					if err != nil {
						mu.Lock()
						errorsMails = append(errorsMails, email)
						error++
						total++
						defer mu.Unlock()
						return
					}
					httpTransport.Dial = dialer.Dial
				}
			}

			client := &http.Client{Transport: httpTransport}

			if modules.IsProxylessMode {
				client = &http.Client{}
			}

			req, _ := http.NewRequest("GET", "https://spclient.wg.spotify.com/signup/public/v1/account?validate=1&email=" + email, nil)

			req.Header.Add("Accept-Language", "en-US,en;q=0.9")
			req.Close = true
			resp, err := client.Do(req)

			if err != nil {
				mu.Lock()
				errorsMails = append(errorsMails, email)
				error++
				total++
				defer mu.Unlock()
				return
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					mu.Lock()
					errorsMails = append(errorsMails, email)
					error++
					total++
					defer mu.Unlock()
					return
				}
			}(resp.Body)

			body, _ := ioutil.ReadAll(resp.Body)
			stringBody := string(body)

			if strings.Contains(stringBody, "That email is already registered to an account") {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()
			} else if strings.Contains(stringBody,"status") {
				mu.Lock()
				bad++
				badLog.Println("[BAD] " + strconv.FormatInt(total, 10) +  " - " + email)
				badFile.Write([]byte(email + "\n"))
				badMails = append(badMails, email)
				defer mu.Unlock()
			} else {
				mu.Lock()
				error++
				errorsMails = append(errorsMails, email)
				defer mu.Unlock()
			}

			total++
			_, _ = utils.SetConsoleTitle("Larez v2.0 | Checked:" + strconv.FormatInt(total, 10) + " - Hits: "+ strconv.FormatInt(good, 10) +" - Bad: "+strconv.FormatInt(bad, 10)+" | " + "Errors: " + strconv.FormatInt(int64(len(errorsMails)), 10))

			if len(errorsMails) > 50 {
				errLog.Println("\nError checker will start in few seconds ! please wait...")
				time.Sleep(10)
				StartEmailChecker(errorsMails, proxies, "recheck")
				utils.ClearConsole()
			}
		},
		emails,
		&wg,
		proxies,
		&mu,
	)
}