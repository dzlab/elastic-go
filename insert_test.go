package elastic

import (
	"testing"
)

// test for insert query
func TestInsert(t *testing.T) {
	actual := []string{
		newInsert().Document(1, Dict{"title": "War and Peace", "author": "Leo Tolstoy"}).String(),
	}
	expected := []string{
		`{"author":"Leo Tolstoy","title":"War and Peace"}`,
	}
	equals(t, actual, expected)
}
