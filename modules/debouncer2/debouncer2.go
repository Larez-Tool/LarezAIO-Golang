package debouncer2

import (
	"awesomeProject/utils"
	"github.com/pterm/pterm"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var Good int64 = 0
var Bad int64 = 0
var Total int64 = 0
var Error int64 = 0

var GoodMails []string
var BadMails []string
var ErrorsMails []string

func StartEmailChecker(email string, proxy string)  {
	goodLog := pterm.NewStyle(pterm.FgGreen)
	badLog := pterm.NewStyle(pterm.FgRed)


	proxyUrl, err := url.Parse("http://" + proxy)

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	req, _ := http.NewRequest("GET", "https://services.postcodeanywhere.co.uk/EmailValidation/Interactive/Validate/v2.00/json3.ws?&Key=ZW43-RN48-NG37-PJ79&Email=" + email + "&Timeout=15000&callback=EmailValidation_Interactive_Validate_v2_00End", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36 Edg/92.0.902.84")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://account.loqate.com/")
	req.Close = true
	resp, err := client.Do(req)

	if err != nil {
		ErrorsMails = append(ErrorsMails, email)
		Error++
		Total++
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ErrorsMails = append(ErrorsMails, email)
			Error++
			Total++
			return
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)

	if strings.Contains(stringBody, "Email address was fully validated") {
		Good++
		goodLog.Println("[GOOD] " + strconv.FormatInt(Total, 10) +  " - " + email)
		GoodMails = append(GoodMails, email)
	} else if strings.Contains(stringBody, "Email Address is not valid") {
		Bad++
		badLog.Println("[BAD] " + strconv.FormatInt(Total, 10) +  " - " + email)
		BadMails = append(BadMails, email)
	} else {
		Error++
		ErrorsMails = append(ErrorsMails, email)
	}

	Total++
	_, _ = utils.SetConsoleTitle("Larez v1.6 | Checked:" + strconv.FormatInt(Total, 10) + " - Hits: "+ strconv.FormatInt(Good, 10) +" - Bad: "+strconv.FormatInt(Bad, 10)+" | " + "Errors: " + strconv.FormatInt(Error, 10))
}
