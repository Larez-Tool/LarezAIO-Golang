package netflix

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

	var debounceIdentifier = "Netflix" + modules.DebounceIdentifier

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

			req, _ := http.NewRequest("GET", "https://www.netflix.com/", nil)
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

			var cookieString = ""
			var cookieNbr = 0

			for _, cookie := range resp.Cookies() {
				if cookieNbr == 0 {
					cookieString = cookieString + cookie.Name + "=" + cookie.Value
					cookieNbr++
				} else {
					cookieString = cookieString + ";" + cookie.Name + "=" + cookie.Value
				}
			}

			body1, _ := ioutil.ReadAll(resp.Body)
			stringBody1 := string(body1)

			authURL := strings.Split(utils.GetStringInBetween(stringBody1, "authURL", `,"isDVD":`), ",")[0]
			authURL = strings.ReplaceAll(strings.ReplaceAll(authURL, `"`, ""), ":", "")
			esn := "NFCDIE-02-" + utils.StringWithCharset(30, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

			if authURL == "" {
				mu.Lock()
				errorsMails = append(errorsMails, email)
				error++
				total++
				defer mu.Unlock()
				return
			}

			authURL, _ = strconv.Unquote(`"` + authURL + `"`)

			payload := strings.NewReader(`{"flow":"signupSimplicity","mode":"welcome","authURL":"` + authURL +`","action":"saveAction","fields":{"email":{"value":"` + email + `"}}}`)

			req2, _ := http.NewRequest("POST", "https://www.netflix.com/api/shakti/v2d4d7c3f/flowendpoint?flow=signupSimplicity&mode=welcome&landingURL=/in/&landingOrigin=https://www.netflix.com&inapp=false&esn=" + esn + "&languages=en-IN&netflixClientPlatform=browser", payload)

			req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")
			req.Header.Add("Accept", "*/*")
			req2.Header.Add("Accept-Language", "en-US,en;q=0.5")
			req2.Header.Add("Content-Type", "application/json")
			req2.Header.Add("X-Netflix.Client.Request.Name", "ui/xhrUnclassified")
			req2.Header.Add("X-Requested-With", "XMLHttpRequest")
			req2.Header.Add("Origin", "https://www.netflix.com")
			req2.Header.Add("Connection", "keep-alive")
			req2.Header.Add("Referer", "https://www.netflix.com/in/")
			req2.Header.Add("Cookie", cookieString)
			req.Header.Add("Accept-Charset", "utf-8")
			req2.Close = true
			resp2, err := client.Do(req2)

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
			}(resp2.Body)

			body2, _ := ioutil.ReadAll(resp2.Body)
			stringBody2 := string(body2)

			if strings.Contains(stringBody2, "switchFlow") {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()
			} else if strings.Contains(stringBody2, "registrationWithContext") {
				mu.Lock()
				bad++
				badLog.Println("[BAD] " + strconv.FormatInt(total, 10) +  " - " + email)
				badMails = append(badMails, email)
				badFile.Write([]byte(email + "\n"))
				defer mu.Unlock()
			} else {
				mu.Lock()
				error++
				errorsMails = append(errorsMails, email)
				defer mu.Unlock()
			}

			total++
			_, _ = utils.SetConsoleTitle("Larez v2.0 | Checked:" + strconv.FormatInt(total, 10) + " - Hits: "+ strconv.FormatInt(good, 10) +" - Bad: "+strconv.FormatInt(bad, 10)+" | " + "Errors: " + strconv.FormatInt(int64(len(errorsMails)), 10))
		},
		emails,
		&wg,
		proxies,
		&mu,
	)

	errLog.Println("Finished Larezed " + strconv.FormatInt(total, 10) + " mails! | Goods: " + strconv.FormatInt(good, 10) + " mails - Bads: " + strconv.FormatInt(bad, 10) + " mails ! | Errors: " + strconv.FormatInt(int64(len(errorsMails)), 10))

	if len(errorsMails) > 50 {
		errLog.Println("\nError checker will start in few seconds ! please wait...")
		time.Sleep(10)
		StartEmailChecker(errorsMails, proxies, "recheck")
		utils.ClearConsole()
	}
}