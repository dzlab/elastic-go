package elastic

import (
	"testing"
)

// test for general queries
func TestGeneral(t *testing.T) {
	// given input
	input := []string{
		String(NewQuery("").KV()),
		String(NewQuery("").Add("argument", "value").KV()),
		String(NewQuery("").AddQuery(NewQuery("query").AddQuery(NewQuery("match_all"))).KV()),
	}
	// expected result
	output := []string{
		"{}",
		`{"argument":"value"}`,
		`{"query":{"match_all":{}}}`,
	}
	for i := 0; i < len(input); i++ {
		if input[i] != output[i] {
			t.Error("Should be equal", input[i], output[i])
		}
	}
}

// test for bool clauses
func TestBool(t *testing.T) {
	input := []string{
		String(NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match").Add("tweet", "elasticsearch"))).KV()),
		String(NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match_all")).AddMustNot(NewQuery("match")).AddShould(NewQuery("match"))).KV()),
	}
	output := []string{
		`{"bool":{"must":{"match":{"tweet":"elasticsearch"}}}}`,
		`{"bool":{"must":{"match_all":{}},"must_not":{"match":{}},"should":{"match":{}}}}`,
	}
	equals(t, input, output)
}

// test for 'term', 'terms' and 'exists' filters
func TestFilters(t *testing.T) {
	actual := []string{
		String(NewQuery("").AddQuery(NewTerm().Add("age", 26)).KV()),
		String(NewQuery("").AddQuery(NewTerms().AddMultiple("tag", "search", "full_text", "nosql")).KV()),
		String(NewQuery("").AddQuery(NewExists().Add("field", "title")).KV()),
	}
	expected := []string{
		`{"term":{"age":26}}`,
		`{"terms":{"tag":["search","full_text","nosql"]}}`,
		`{"exists":{"field":"title"}}`,
	}
	equals(t, actual, expected)
}
