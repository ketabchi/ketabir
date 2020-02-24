package melli

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/ketabchi/ketabir/api"
)

type Book struct {
	url string
	doc *goquery.Document
}

var NoBookErr = errors.New("no book with this isbn")

func NewBookByISBN(isbn string) (*Book, error) {
	url, err := api.GetBookURLByISBN(isbn)
	if err != nil {
		return nil, err
	}
	if url == "" {
		return nil, NoBookErr
	}

	return NewBook(url)
}

func NewBook(url string) (*Book, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	return &Book{url: url, doc: doc}, nil
}

func (b *Book) Name() string {
	return b.doc.Find("#ctl00_ContentPlaceHolder1_lblBookTitle").Text()
}

func (b *Book) Publisher() (publisher string) {
	return b.doc.Find("#ctl00_ContentPlaceHolder1_rptPublisher_ctl00_NameLabel").Text()
}

func (b *Book) Authors() []string {
	authors := make([]string, 0)
	b.doc.Find("#ctl00_ContentPlaceHolder1_rptAuthor span").EachWithBreak(
		func(i int, sel *goquery.Selection) bool {
			s := sel.Text()
			if !strings.Contains(s, "نويسنده:") {
				return true
			}

			s = strings.Replace(s, "نويسنده:", "", -1)
			authors = append(authors, strings.TrimSpace(s))

			return true
		})

	return authors
}

func (b *Book) Translators() []string {
	translators := make([]string, 0)
	b.doc.Find("#ctl00_ContentPlaceHolder1_rptAuthor span").EachWithBreak(
		func(i int, sel *goquery.Selection) bool {
			s := sel.Text()
			if !strings.Contains(s, "مترجم:") {
				return true
			}

			s = strings.Replace(s, "مترجم:", "", -1)
			translators = append(translators, strings.TrimSpace(s))

			return true
		})

	return translators
}

func (b *Book) Link() string {
	return b.url
}
