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
	expected1 := []string{
		"?search_type=scan&scroll=1m",
	}
	expected2 := []string{
		"?scroll=1m&search_type=scan",
	}
	for i := 0; i < len(actual); i++ {
		if !(actual[i] == expected1[i] || actual[i] == expected2[i]) {
			t.Errorf("%s Should be equal\n%s or %s", actual[i], expected1[i], expected2[i])
		}
	}
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

// test for query clauses
func TestQuery(t *testing.T) {
	actual := []string{
		newQuery().AddQuery(NewQuery("dis_max").AddQueries("queries", NewQuery("match").Add("title", "Brown fox"), NewQuery("match").Add("body", "Brown fox"))).String(),
	}
	expected := []string{
		`{"dis_max":{"queries":[{"match":{"title":"Brown fox"}},{"match":{"body":"Brown fox"}}]}}`,
	}

	equals(t, actual, expected)
}

// test for bool clauses
func TestBool(t *testing.T) {
	input := []string{
		NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match").Add("tweet", "elasticsearch"))).String(),
		NewQuery("").AddQuery(NewBool().AddMust(NewQuery("match_all")).AddMustNot(NewQuery("match")).AddShould(NewQuery("match"))).String(),
		NewQuery("").AddQuery(NewBool().AddShould(NewTerm().Add("price", 20)).AddShould(NewTerm().Add("productID", "XHDK-A-1293-#fJ3")).AddShould(NewTerm().Add("category", "smartphone"))).String(),
	}
	output := []string{
		`{"bool":{"must":{"match":{"tweet":"elasticsearch"}}}}`,
		`{"bool":{"must":{"match_all":{}},"must_not":{"match":{}},"should":{"match":{}}}}`,
		`{"bool":{"should":[{"term":{"price":20}},{"term":{"productID":"XHDK-A-1293-#fJ3"}},{"term":{"category":"smartphone"}}]}}`,
	}
	equals(t, input, output)
}

// test for 'term', 'terms' and 'exists' filters
func TestFilters(t *testing.T) {
	actual := []string{
		NewQuery("").AddQuery(NewTerm().Add("age", 26)).String(),
		NewQuery("").AddQuery(NewTerms().AddMultiple("tag", "search", "full_text", "nosql")).String(),
		NewQuery("").AddQuery(NewTerms().AddMultiple("price", 20, 30)).String(),
		NewQuery("").AddQuery(NewExists().Add("field", "title")).String(),
	}
	expected := []string{
		`{"term":{"age":26}}`,
		`{"terms":{"tag":["search","full_text","nosql"]}}`,
		`{"terms":{"price":[20,30]}}`,
		`{"exists":{"field":"title"}}`,
	}
	equals(t, actual, expected)
}
