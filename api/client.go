package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{}

func SetClient(c *http.Client) {
	client = c
}

var chapRe = regexp.MustCompile(`چاپ ([\d]+) سال`)

func GetBookURLByISBN(isbn string) (string, error) {
	if isbn == "" {
		return "", nil
	}

	body, err := createPostBodyISBN(isbn)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://ketab.ir/Search.aspx", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("not 200 on sending request")
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}

	el := doc.Find("#ctl00_ContentPlaceHolder1_DataList1 > tbody > tr")
	if el.Length() == 0 {
		return "", nil
	}

	link, chap := "", 0
	el.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		s := sel.Find("tr:nth-child(2) span").Text()
		ss := chapRe.FindStringSubmatch(s)
		if len(ss) < 2 {
			return true
		}

		ch, _ := strconv.Atoi(ss[1])
		if ch > chap {
			chap = ch
			href, _ := sel.Find(".HyperLink2").Attr("href")

			link = fmt.Sprintf("http://ketab.ir%s", href)
		}

		return true
	})

	return link, nil
}

func createPostBodyISBN(isbn string) (io.Reader, error) {
	res, err := client.Get("http://ketab.ir/Search.aspx")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = setClientCookies(res.Cookies()); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	v1, _ := doc.Find("#__VIEWSTATE").Attr("value")
	v2, _ := doc.Find("#__VIEWSTATEGENERATOR").Attr("value")
	v3, _ := doc.Find("#__EVENTVALIDATION").Attr("value")

	data := url.Values{}
	data.Set("__VIEWSTATE", v1)
	data.Set("__VIEWSTATEGENERATOR", v2)
	data.Set("__EVENTVALIDATION", v3)
	data.Set("ctl00$SiteMenu$Search$DropDownFieldList", "1")
	data.Set("ctl00$SiteMenu$Search$DDLTypeSearch", "1")
	data.Set("ctl00$ContentPlaceHolder1$TxtIsbn", isbn)
	data.Set("ctl00$ContentPlaceHolder1$drpDewey", "-1")
	data.Set("ctl00$ContentPlaceHolder1$drpFromIssueYear", "57")
	data.Set("ctl00$ContentPlaceHolder1$drpFromIssueMonth", "01")
	data.Set("ctl00$ContentPlaceHolder1$drpFromIssueDay", "01")
	data.Set("ctl00$ContentPlaceHolder1$drpToIssueYear", "98")
	data.Set("ctl00$ContentPlaceHolder1$drpToIssueMonth", "12")
	data.Set("ctl00$ContentPlaceHolder1$drpToIssueDay", "30")
	data.Set("ctl00$ContentPlaceHolder1$drLanguage", "0")
	data.Set("ctl00$ContentPlaceHolder1$DrSort", "1")
	data.Set("ctl00$ContentPlaceHolder1$DrPageSize", "100")
	data.Set("ctl00$ContentPlaceHolder1$BtnSearch", "")

	return strings.NewReader(data.Encode()), nil
}

func setClientCookies(cs []*http.Cookie) error {
	u, err := url.Parse("http://ketab.ir")
	if err != nil {
		return err
	}

	cj, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	cj.SetCookies(u, cs)
	client.Jar = cj

	return nil
}
