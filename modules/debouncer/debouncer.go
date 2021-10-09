package debouncer

import (
	"awesomeProject/modules"
	"awesomeProject/utils"
	"fmt"
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

	var debounceIdentifier = "Debounce" + modules.DebounceIdentifier

	_, err := os.Stat("./results/good/" + debounceIdentifier)
	_, err = os.Stat("./results/bad/" + debounceIdentifier)
	var goodFile *os.File
	var badFile *os.File

	if os.IsNotExist(err) {
		goodFile, _ = os.Create("./results/good/" + debounceIdentifier)
		badFile, _ = os.Create("./results/bad/" + debounceIdentifier)
		fmt.Println("Dir doesn't exist, waiting..")
	} else {
		goodFile, _ = os.Open("./results/good/" + debounceIdentifier)
		badFile, _ = os.Open("./results/bad/" + debounceIdentifier)
		fmt.Println("Dir exist, waiting..")
	}

	var goodMails []string
	var badMails []string
	var errorsMails []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	modules.EmailCheck(
		func (email string, wg *sync.WaitGroup, proxyS string, mu *sync.Mutex) {
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

			payload := strings.NewReader(`Email=` + email + `&Key=` + modules.DebounceKey)

			req, _ := http.NewRequest("POST", "https://services.postcodeanywhere.co.uk/EmailValidation/Interactive/Validate/v2.00/json3.ws", payload)
			req.Close = true

			req.Header.Add("accept", "*/*")
			req.Header.Set("Accept-Encoding", "identity")
			req.Header.Add("accept-language", "fr,fr-FR;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
			req.Header.Add("content-length", "52")
			req.Header.Add("content-type", "application/x-www-form-urlencoded")
			req.Header.Add("origin", "https://www.loqate.com")
			req.Header.Add("referer", "https://www.loqate.com/")
			req.Header.Add("sec-ch-ua", "\"Microsoft Edge\";v=\"93\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"93\"")
			req.Header.Add("sec-ch-ua-mobile", "?0")
			req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
			req.Header.Add("sec-fetch-dest", "empty")
			req.Header.Add("sec-fetch-mode", "cors")
			req.Header.Add("sec-fetch-site", "cross-site")
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47")
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

			if strings.Contains(stringBody, "Email address was fully validated") {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()

			} else if strings.Contains(stringBody, "Email Address is not valid") {
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