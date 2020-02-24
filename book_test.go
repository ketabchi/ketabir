package ketabir

import (
	"testing"

	"github.com/ketabchi/util"
)

func TestName(t *testing.T) {
	tests := []struct {
		url string
		exp string
	}{
		{
			"http://ketab.ir/bookview.aspx?bookid=2273622",
			"سمفونی مردگان",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2425864",
			"شغل مناسب شما: با توجه به ویژگی‌های شخصیتی خود کارتان را انتخاب کنید، جویندگان کار چگونه کار مورد علاقه خود را انتخاب کنند",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2402503",
			"مدیریت اجرایی MBA) for dummies)",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2359448",
			"شدن",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=1871774",
			"طلبه زیستن: پژوهشی مقدماتی در سنخ‌شناسی جامعه‌شناختی زیست‌طلبگی",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=1809369",
			"ارتباط رو در رو: کلید موفقیت برای مدیریت موثر و کارا (مجموعه مقالاتی از دانشگاه هاروارد)",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=1911364",
			"دریدا و فلسفه",
		},
	}

	for i, test := range tests {
		book, err := NewBook(test.url)
		if err != nil {
			t.Errorf("Test %d: Error on creating book from %s: %s",
				i, test.url, err)
		}
		if name := book.Name(); name != test.exp {
			t.Errorf("Test %d: Expected book name '%s', but got '%s'",
				i, test.exp, name)
			t.Logf("\n%q\n%q", test.exp, name)
		}
	}
}

func TestPublisher(t *testing.T) {
	tests := []struct {
		url string
		exp string
	}{
		{
			"http://ketab.ir/bookview.aspx?bookid=2273622",
			"ققنوس",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2425864",
			"نقش و نگار",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2402503",
			"آوند دانش",
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2359448",
			"مهراندیش",
		},
	}

	for i, test := range tests {
		book, err := NewBook(test.url)
		if err != nil {
			t.Errorf("Test %d: Error on creating book from %s: %s",
				i, test.url, err)
		}
		if name := book.Publisher(); name != test.exp {
			t.Errorf("Test %d: Expected publisher name '%s', but got '%s'",
				i, test.exp, name)
			t.Logf("\n%q\n%q", test.exp, name)
		}
	}
}

func TestAuthors(t *testing.T) {
	tests := []struct {
		url string
		exp []string
	}{
		{
			"http://ketab.ir/bookview.aspx?bookid=2327303",
			[]string{"گری نورتفیلد"},
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2402503",
			[]string{"کتلین‌آر. الن", "پیتر اکونومی"},
		},
	}

	for i, test := range tests {
		book, err := NewBook(test.url)
		if err != nil {
			t.Errorf("Test %d: Error on creating book from %s: %s",
				i, test.url, err)
		}

		authors, ss := make([]string, 0), book.Authors()
		for _, s := range ss {
			authors = append(authors, s[0])
		}

		if !util.CheckSliceEq(authors, test.exp) {
			t.Errorf("Test %d: Expected authors %q, but got %q",
				i, test.exp, authors)
		}
	}
}

func TestTranslators(t *testing.T) {
	tests := []struct {
		url string
		exp []string
	}{
		{
			"http://ketab.ir/bookview.aspx?bookid=1839057",
			[]string{"محمدرضا طبیب‌زاده"},
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2368628",
			[]string{},
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2319963",
			[]string{"پریسا صیادی", "سرور صیادی"},
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2364768",
			[]string{"عادل فردوسی‌پور", "علی شهروزستوده", "بهزاد توکلی‌نیشابوری"},
		},
		{
			"http://ketab.ir/bookview.aspx?bookid=2313586",
			[]string{"امیرحسین میرزائیان", "عبدالرضا شهبازی"},
		},
	}

	for i, test := range tests {
		book, err := NewBook(test.url)
		if err != nil {
			t.Errorf("Test %d: Error on creating book from %s: %s",
				i, test.url, err)
		}
		if translators := book.Translators(); !util.CheckSliceEq(translators, test.exp) {
			t.Errorf("Test %d: Expected translators %q, but got %q",
				i, test.exp, translators)
		}
	}
}
