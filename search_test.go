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
		"{\"argument\":\"value\"}",
		"{\"query\":{\"match_all\":{}}}",
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
		"{\"bool\":{\"must\":{\"match\":{\"tweet\":\"elasticsearch\"}}}}",
		"{\"bool\":{\"must\":{\"match_all\":{}},\"must_not\":{\"match\":{}},\"should\":{\"match\":{}}}}",
	}
	for i := 0; i < len(input); i++ {
		if input[i] != output[i] {
			t.Error("Should be equal", input[i], output[i])
		}
	}
}
