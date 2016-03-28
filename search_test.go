package elastic

import (
	"testing"
)

// test for general queries
func TestGeneral(t *testing.T) {
	// given input
	input := []string{
		NewQuery("").String(),
		NewQuery("").Add("argument", "value").String(),
		NewQuery("").AddQuery(NewQuery("query").AddQuery(NewQuery("match_all"))).String(),
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

// test for search url
func TestSearchUrl(t *testing.T) {
	actual := []string{
		newSearch().AddParam(SEARCH_TYPE, "scan").AddParam(SCROLL, "1m").urlString(),
	}
	expected := []string{
		"?search_type=scan&scroll=1m",
	}
	equals(t, actual, expected)
}

// test for search queries
func TestSearch(t *testing.T) {
	actual := []string{
		newSearch().AddParam("search_type", "scan").AddParam("scroll", "1m").AddQuery(NewQuery("query").AddQuery(NewQuery("range").AddQuery(NewQuery("data").Add("gte", "2014-01-01").Add("lt", "2014-02-01")))).Add("size", 1000).String(),
		newSearch().AddQuery(NewQuery("query").AddQuery(NewQuery("match_all"))).AddSource("title").AddSource("created").String(),
	}
	expected := []string{
		`{"query":{"range":{"data":{"gte":"2014-01-01","lt":"2014-02-01"}}},"size":1000}`,
		`{"_source":["title","created"],"query":{"match_all":{}}}`,
	}
	equals(t, actual, expected)
}

// test for bool clauses
func TestBool(t *testing.T) {
	input := []string{
		NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match").Add("tweet", "elasticsearch"))).String(),
		NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match_all")).AddMustNot(NewQuery("match")).AddShould(NewQuery("match"))).String(),
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
