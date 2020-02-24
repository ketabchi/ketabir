package api

import "testing"

func TestGetBookURLByISBN(t *testing.T) {
	tests := []struct {
		isbn string
		exp  string
	}{
		{
			"9789643113445",
			"http://ketab.ir/bookview.aspx?bookid=2453934",
		},
		{
			"9789646235793",
			"http://ketab.ir/bookview.aspx?bookid=2425864",
		},
	}

	for i, test := range tests {
		url, err := GetBookURLByISBN(test.isbn)
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
