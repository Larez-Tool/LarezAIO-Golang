package numdebouncer

import (
	"awesomeProject/utils"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var good int64 = 0
var bad int64 = 0
var total int64 = 0
var error int64 = 0

func StartNumDebouncerNoGoroutine(emails []string, proxies []string)  {
	goodLog := pterm.NewStyle(pterm.FgGreen)
	badLog := pterm.NewStyle(pterm.FgRed)
	errLog := pterm.NewStyle(pterm.FgLightYellow)
	//startingDebounceTime := time.Now().Unix()

	var goodMails []string
	var badMails []string
	var errorsMails []string
	proxyOffset := 0
	proxyOffsetMax := len(proxies) - 1

	for _, email := range emails {
		if proxyOffset == proxyOffsetMax  {
			proxyOffset = 0
		} else {
			proxyOffset++
		}
		fmt.Println(proxies[proxyOffset])
		proxyUrl, err := url.Parse("http://" + proxies[proxyOffset])

		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

		payload := strings.NewReader(`phoneNumber=` + email + `&isoCode=FR&version=3`)
		req, err := http.NewRequest("POST", "https://www.edq.com/phonedemo/validatephonenumber", payload)

		req.Header.Set("Content-Length", "45")
		req.Header.Set("Requestverificationtoken", "r1eHeSlclrnjhe7gXxdXMmrFgR2iziUO5scAi-Km6QYX51mFp7igGHuatrlct46NZoB0ccSFNnTpk63SLj21_79JmZU1")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("X-Edq-Site-Instance", "Experian.GlobalCms.Website-Production-NorthEurope/5.21.1")
		req.Header.Set("Origin", "https://www.edq.com")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Referer", "https://www.edq.com/phone-verification/bulk-phone-verification/")
		req.Header.Set("Cookie", "edqcmssession=nchha4iz301iiwx53dmkce03; edqcmsuser=marketingId=80191aaf-3c7c-4cf4-9348-d202548892c2; edqcmsxsrf=wM5EXvFLSain3pt_TV4aefKhamLSc7eDNoIpvjQxgdgdvYnOswZaXpBvDWTD284YfQpVPEKBDfoKVSa60tEICWwLZQc1; ARRAffinity=f56d556fd4dfe69b571fc2d95d4031d06b65aed5e94c05b3316852e9822eb105; ARRAffinitySameSite=f56d556fd4dfe69b571fc2d95d4031d06b65aed5e94c05b3316852e9822eb105; visid_incap_1333723=yJFkJdqJR82nQfb3L8f+ojIvNmEAAAAAQUIPAAAAAADV7P1M6tYqQzxuP2fITC1Y; nlbi_1333723=2o6lPur3oW0IKnhf7yko6gAAAACghoaTs7KjexZMT9FA/wcI; incap_ses_189_1333723=gU1acclp9AzmhSNp/XafAjQvNmEAAAAAO06knrw8rbBwL1icTguQhQ==; incap_ses_1516_2108539=mGVBaTUuZSixnwd1A+oJFTUvNmEAAAAAaToXw9N54mfnHDCc0vE++Q==; visid_incap_2108539=65ub3ZFaQN+dVsTc5p2IBjQvNmEAAAAAQUIPAAAAAAARkAOctYYRe7FgHDnzFR7q; incap_ses_189_2108539=r1dnUuxPKgpThyNp/XafAjUvNmEAAAAAiyrtvjqd2gFaHfk5U4uOcw==; _ga=GA1.2.1518847641.1630940983; _gid=GA1.2.473225663.1630940983; _gcl_au=1.1.1806162465.1630940983; googleAnalyticsClientId=1518847641.1630940983; bf_lead=1vunoi226roo00; _gd_visitor=5bed5ea9-97ee-48d5-8851-948c45218bf4; _gd_session=2ceadfd3-e212-4551-85fc-29f3990c06d1; calltrk_referrer=https%3A//www.google.com/; calltrk_landing=https%3A//www.edq.com/email-verification/; _hjid=44d7d00b-e6d7-42f7-b571-3d56bdc8b988; _hjFirstSeen=1; _gd_svisitor=670fdd58301a0000372f3661a20100000b1d0900; _hjIncludedInPageviewSample=1; _hjAbsoluteSessionInProgress=0; _hjIncludedInSessionSample=1; _an_uid=0; bf_visit=9bgpdh4n2rc00; incap_ses_1516_1333723=OXg8cZui3hhLHQh1A+oJFccvNmEAAAAAj8QlpR2ZtdlrpEPKZZgFTg==; _dc_gtm_UA-25829750-3=1; _uetsid=773f51600f2411ec85404d3fe4dde24d; _uetvid=773f9e900f2411eca8c7edece5b35deb; _gali=phone-validate-search")
		req.Header.Add("Accept-Charset", "utf-8")
		req.Header.Add("Accept", "*/*")
		req.Close = true
		resp, err := client.Do(req)

		if err != nil {
			errLog.Println("[ERROR] " + strconv.FormatInt(total, 10) +  " - " + email + " TYPE: Invalid response 1")
			errorsMails = append(errorsMails, email)
			error++
			total++
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				errLog.Println("[ERROR] " + strconv.FormatInt(total, 10) +  " - " + email + " TYPE: Invalid response 2")
				errorsMails = append(errorsMails, email)
				error++
				total++
				return
			}
		}(resp.Body)

		body, _ := ioutil.ReadAll(resp.Body)
		stringBody := string(body)

		fmt.Println(stringBody)

		if strings.Contains(stringBody, "verified") {
			good++
			goodLog.Println("[GOOD] " + strconv.FormatInt(total, 10) +  " - " + email)
			goodMails = append(goodMails, email)
		} else if strings.Contains(stringBody, "unknown") || strings.Contains(stringBody, "dead") || strings.Contains(stringBody, "invalid"){
			bad++
			badLog.Println("[BAD] " + strconv.FormatInt(total, 10) +  " - " + email)
			badMails = append(badMails, email)
		} else {
			errLog.Println("[ERROR] " + strconv.FormatInt(total, 10) +  " - " + email + " TYPE: Invalid response 3")
			error++
			errorsMails = append(errorsMails, email)
		}

		total++
		_, _ = utils.SetConsoleTitle("Larez v1.6 | Checked:" + strconv.FormatInt(total, 10) + " - Hits: "+ strconv.FormatInt(good, 10) +" - Bad: "+strconv.FormatInt(bad, 10)+" | " + "Errors: " + strconv.FormatInt(error, 10))
	}
}

