package disney

import (
	"awesomeProject/modules"
	"awesomeProject/utils"
	"encoding/json"
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

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func StartEmailChecker(emails []string, proxies []string, mode string)  {
	goodLog := pterm.NewStyle(pterm.FgGreen)
	badLog := pterm.NewStyle(pterm.FgRed)
	errLog := pterm.NewStyle(pterm.FgLightYellow)

	var debounceIdentifier = "Disney" + modules.DebounceIdentifier

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

			req, _ := http.NewRequest("GET", "https://www.disneyplus.com/fr-fr/login", nil)
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

			clientApiKey := strings.Split(utils.GetStringInBetween(stringBody1, "clientApiKey", `,"environment":`), ",")[0]
			clientApiKey = strings.ReplaceAll(strings.ReplaceAll(clientApiKey, `"`, ""), ":", "")

			if clientApiKey == "" {
				mu.Lock()
				errorsMails = append(errorsMails, email)
				error++
				total++
				defer mu.Unlock()
				return
			}

			payload := strings.NewReader("grant_type=refresh_token&latitude=0&longitude=0&platform=browser&refresh_token=eyJ6aXAiOiJERUYiLCJraWQiOiJLcTYtNW1Ia3BxOXdzLUtsSUUyaGJHYkRIZFduRjU3UjZHY1h6aFlvZi04IiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..vmczy20GX6caChWT.9DqAmfeoJAui2vppOLMOjaCljPQeFm0Il7mX9uhbJoWQAk_sbcFlESqOeH5v7w46rL5suCHJJn0SJiAjJFTVHGTiLsI5OC7ii2qNrRrgZ9mE-2fc535zoe3jgKeDCYdXSrofo5MXO9jtb9VtqWcfE9TrYiINstjAsRspt4EKTFjUe1z6betWexVwh-mgNPrrvehQPeZoPOVf41tzBzzqYQg_2rroBCwdAPylLC1x8qZPmI4fqyWX0np9VNRrr2S1rpyjzWxXJvsoLDgB4QU1JJMm0LOaLr6DYK8tQbpyr5kvqqyR2iqqqPTou-v6qdjtNNawNnwyuNXIEEB6qNEQDb7MgvLzscrZngT-uOa8IW6uI6LnoL0c0XI_EqB-DbFo2vE22c3DVbmtFQitdx9BPoxlCHHtzyjmZZRn44HQz5mhFgruoRCQ3naI2yKFenlzSr7swWUOy0jwNePjrdBp0zr0YEKH0ldWdrYHBUl6GLEqL2yi6krMLXI9-diZVdj_3BYu3mjn58EtJ5Pat7Ywun5uJyNHFJMP0D-36xWNjNZB6OQrfz2q25vgWR8vFS2OmtW8v5EzxInIqc_wvznmmvfKJhRrzW0HKccWtZUUgLMh1A302VxCoaRc-O7s_xKhqGrLqRPL5BGR-mEv19s_pgX9Mn2zC6vQE-PLfpNSNaTSr6SrCq1FCGw2_aGFv9Qvx_8aHnLbrViIZfJpFNHMIsiGsnEZBBC8l9JCNtW8uJexeAX8u10ZMLCafMwmnmfr8lZZAeTfe3MBGZpi-bxrr_Sy1Hdtui-fcf1yJIUPocPQ46fxRFFjvmPDF1MlY-HvVROIr5qlph7W0HgvzZCkkjDjuVyWp6TSOfAgKATZl82Aj0CVWMB2p8YXD-VN60Uk5b3l4M4-KNWj8zSsuqdOEE1fXShoW1LQZMqW8-y1UsLv44FcTI1A7kWdqWyGlWbdZN5Sw8Y1_n8eCtmQdCMggR6jxJjJI_sl1F53oQYXOtfnlzjJNOhNj-VCYX_bjsOn_Rzbj5aM5UXEn_toCjvDAOkD5v3OCFox0yeOIuJt27w0aeFSKZg9_zWWsuGb5TjbsIm_947EXgrcnoriwn02yVIMcw5sGGj-K0-forSJJg0BoT9fPveiR8XqhngIGFNjmbNRUf6IP8mb6zmgWMju0VeLKQTg-6ApqD9pDH09xVBAJCPBsYlU_IQkeMLm0_6JDTKIJePoGIws7FJ7x_796VqBFeJ4P4CrSUdd310JAoZgMIExAcqzsrEhHbFDLq-ETc2ffkhqly6uIYZq0cD0dhBkhsm_S2svAzhy15bXxOoJpSQhU0rV6bI7Ecs4FrtlgQbau13xWchLHnrdhiramFwIWuVVg3ibtPSzKY17HHObU4XX3AN0DV2EC_5V7-rCajT_oaa1t9yQMi4Q_RfCDFY5WOriwLkZdrA7C5Po3lb8HbaDh30aJqYzR5cIIYJROxzLALPWRH9Fm_J-B7BuWHYqk7OXz7H4iGSFV7Ym30u8swQXZuq9LWvv807Yal_Qg9u2vOJW1tg35eNgDpvlQO84ssj2EUgM42ixkd-l9mJGt3051Y_C-_vbuQRW_rZ8KuIC0G5ZCFv7t0wt9csdEkW_wEKdzp1MvTN82YaHJB88j0pmS-UzrDPT3576U8QbX_oMJMWLcKSXSU4DZW4HjFbvDYzRNFKV-uENNzKNfVjpakSOu9fmCljhhbDbJXqJHTfp8-M7L0fuEXQ1x7HH1Jq1SdnEEAieWe__vbpKmAKSDm5EFy7PkXcg_ZAfdU2rUXOZEuKF5Y6kthmuLgwzsr4EIbT63QWxGtp5dswM5QOV8oLog_8hFm2z93Lwckq9aymRHhDz9xmBaP8upQIs3iY9IpTxf1mrNWM1NGghKFVcdddioz10f7UXo0FmX3n5rA4sIs3pmq0IV7MXFUzPVD9JdI_TEoqY8kYT9wdcEyrm-m6jgllfFf2kcHdwqh18baLvy9MpYoV8il5mzs2ylPS9orDWzAJ2HwcDNxYbr90zCebCjzKr14qUOIPNSwx_ZE8xVXpwyHy_6nL6ovwWCTpLlbwKeK1K8yTm8Z_3iYD8ocOzYf6Nm0LkWXs_WLQjUVNSh5FtESafbP13_tXkDQwfllw0lA-Q3BVOgbCZbxEnE84wMRFMxniAOTi57mwVVMibQFiGEFHAE3H_TlvrEk3gkehuLRntNOKlA5Q1nvxa7HC-8FOPQpO4C3duM5nloWmqGw4KawPuOqHtM1MLJXWxZPF4vslF8u8oKwqqc16pJDLrUuTX2JTcLHFnDBjq1VDReXwBbcljKWnT9rJ5l1tR_yplc9cP9HvMdk0yDkpoFK30RWmLilOZj2kJRHjbsWIN_FXtnpR0wRXXwzaVmS8MsX3yz8PBj0lmnGSvFfPScatjA6V2XLvHgJcFjgb2imyY_PPjMX6XdxXrWDyWDt9Ty-29QGYA-cdRIj4bhDll_L2I6zVLPyHSBrRk08PfQ_ujJNUjN68badjjW_HPtJm85GRYjISNMl4yrWSW48FGPkKbvvdmhufb64xIEw6NLgkMvZsG0R0yiSQ.Ei_hPy3JrhWFAY13iJCY2w")

			req2, _ := http.NewRequest("POST", "https://global.edge.bamgrid.com/token" , payload)

			req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0")
			req2.Header.Add("x-bamsdk-platform", "windows")
			req2.Header.Add("x-bamsdk-client-id", "disney-svod-3d9324fc")
			req2.Header.Add("x-application-version", "1.1.2")
			req2.Header.Add("sec-ch-ua-mobile", "?0")
			req2.Header.Add("authorization", "Bearer " + clientApiKey)
			req2.Header.Add("content-type", "application/x-www-form-urlencoded")
			req2.Header.Add("x-bamsdk-version", "7.0")
			req2.Header.Add("accept", "application/json")
			req2.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54")
			req2.Header.Add("x-dss-edge-accept", "vnd.dss.edge+json; version=2")

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

			var tRes tokenResponse
			json.Unmarshal([]byte(stringBody2), &tRes)

			accessToken := tRes.AccessToken


			payload2 := strings.NewReader(`{"query":"query Check($email: String!) {check(email: $email) {operations}}","variables":{"email":"` + email + `"}}`)
			req3, _ := http.NewRequest("POST", "https://global.edge.bamgrid.com/v1/public/graphql", payload2)

			req3.Header.Add("x-bamsdk-platform", "windows")
			req3.Header.Add("x-bamsdk-client-id", "disney-svod-3d9324fc")
			req3.Header.Add("x-application-version", "1.1.2")
			req3.Header.Add("sec-ch-ua-mobile", "?0")
			req3.Header.Add("authorization", accessToken)
			req3.Header.Add("x-bamsdk-platform-id", "browser")
			req3.Header.Add("content-type", "application/json")
			req3.Header.Add("x-bamsdk-version", "7.0")
			req3.Header.Add("accept", "application/json")
			req3.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54")
			req3.Header.Add("x-dss-edge-accept", "vnd.dss.edge+json; version=2")

			req3.Close = true
			resp3, err := client.Do(req3)

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
			}(resp3.Body)

			body3, _ := ioutil.ReadAll(resp3.Body)
			stringBody3 := string(body3)

			if strings.Contains(stringBody3, "Login") || strings.Contains(stringBody3, "OTP") {
				mu.Lock()
				good++
				goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
				goodMails = append(goodMails, email)
				goodFile.Write([]byte(email + "\n"))
				defer mu.Unlock()
			} else if strings.Contains(stringBody3, "Register") {
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
