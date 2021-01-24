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
	"github.com/antzucaro/matchr"
)

type urlFinder struct {
	isbn   string
	title  string
	client http.Client
	doc    *goquery.Document
}

var Client = &http.Client{}

var chapRe = regexp.MustCompile(`چاپ ([\d]+) سال`)

func GetBookURLByISBN(isbn string, args ...string) (string, error) {
	if isbn == "" {
		return "", nil
	}

	uf := &urlFinder{isbn: isbn, client: *Client}
	if len(args) > 0 {
		uf.title = args[0]
	}

	body, err := uf.createBody()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://db.ketab.ir/Search.aspx", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := uf.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("not 200 on sending request")
	}

	uf.doc, err = goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}

	return uf.find(), nil
}

func (uf *urlFinder) createBody() (io.Reader, error) {
	res, err := uf.client.Get("https://db.ketab.ir/Search.aspx")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err = uf.setCookies(res.Cookies()); err != nil {
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
	data.Set("ctl00$ContentPlaceHolder1$TxtIsbn", uf.isbn)
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

func (uf *urlFinder) setCookies(cs []*http.Cookie) error {
	u, err := url.Parse("https://db.ketab.ir")
	if err != nil {
		return err
	}

	cj, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	cj.SetCookies(u, cs)
	uf.client.Jar = cj

	return nil
}

func (uf *urlFinder) find() (link string) {
	chap := 0
	el := uf.doc.Find("#ctl00_ContentPlaceHolder1_DataList1 > tbody > tr")
	el.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		s := sel.Find("tr:nth-child(2) span").Text()
		ss := chapRe.FindStringSubmatch(s)
		if len(ss) < 2 {
			return true
		}

		pass := true
		if uf.title != "" {
			title := sel.Find(".HyperLink2").Text()
			title = strings.TrimSpace(title)
			tmp := matchr.SmithWaterman(uf.title, title)
			tmp /= float64(len([]rune(uf.title)))

			if tmp <= 0.2 {
				pass = false
			}
		}

		if pass {
			ch, _ := strconv.Atoi(ss[1])
			if ch > chap {
				chap = ch
				href, _ := sel.Find(".HyperLink2").Attr("href")

				link = fmt.Sprintf("https://db.ketab.ir%s", href)
			}
		}

		return true
	})
	return
}
