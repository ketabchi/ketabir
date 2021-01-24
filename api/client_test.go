package api

import "testing"

func TestGetBookURLByISBN(t *testing.T) {
	tests := []struct {
		isbn string
		args []string
		exp  string
	}{
		{
			"9789643113445",
			nil,
			"https://db.ketab.ir/bookview.aspx?bookid=2476393",
		},
		{
			"9789646235793",
			nil,
			"https://db.ketab.ir/bookview.aspx?bookid=2425864",
		},
		{
			"9786002571755",
			[]string{"توانبخشی مبتنی بر جامعه شهری (براساس تجربه مددکاران اجتماعی در شهرستان قدس)"},
			"https://db.ketab.ir/bookview.aspx?bookid=2202390",
		},
	}

	for i, test := range tests {
		url, err := GetBookURLByISBN(test.isbn, test.args...)
		if err != nil {
			t.Errorf("Test %d: Error on getting book url by %s isbn: %s",
				i, test.isbn, err)
		}
		if url != test.exp {
			t.Errorf("Test %d: Expected %s but got %s for %s isbn.",
				i, test.exp, url, test.isbn)
		}
	}
}
