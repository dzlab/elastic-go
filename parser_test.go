package elastic

import (
	"testing"
)

// test for search result parser
func TestSearchResultParser(t *testing.T) {
	parser := &SearchResultParser{}
	input := []string{
		`{"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"failed":0},"hits":{"total":1,"max_score":0.50741017,"hits":[{"_index":"my_index","_type":"my_type","_id":"1","_score":0.50741017,"_source":{"name":"Brown foxes"}}]}}`,
		`{"took":1,"timed_out":false,"_shards":{"total":5,"successful":5,"failed":0},"hits":{"total":0,"max_score":null,"hits":[]}}`,
	}
	expected := []interface{}{
		SearchResult{Took: 3, TimedOut: false, Shards: Shard{Total: 1, Successful: 1, Failed: 0}, Hits: Hits{Total: 1, MaxScore: 0.50741017, Hits: []SearchHits{SearchHits{Index: "my_index", Type: "my_type", ID: "1", Score: 0.50741017, Source: Dict{"name": "Brown foxes"}}}}},
		SearchResult{Took: 1, TimedOut: false, Shards: Shard{Total: 5, Successful: 5, Failed: 0}, Hits: Hits{Total: 0, MaxScore: nil, Hits: make([]SearchHits, 0)}},
	}
	checkParsingResult(t, input, parser, expected)
}

// TestInsertResultParser tests for InsertResultParser
func TestInsertResultParser(t *testing.T) {
	parser := &InsertResultParser{}
	// input strings to parse
	input := []string{
		`{"_index":"blogposts","_type":"post","_id":"1","_version":2,"_shards":{"total":2,"successful":1,"failed":0},"created":false}`,
		`{"_index":"my_index","_type":"groups","_id":"1","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"created":true}`,
		`{"acknowledged":true}`,
	}
	// expected results
	expected := []interface{}{
		InsertResult{Index: "blogposts", Doctype: "post", ID: "1", Version: 2, Shards: Shard{Total: 2, Successful: 1, Failed: 0}, Created: false},
		InsertResult{Index: "my_index", Doctype: "groups", ID: "1", Version: 1, Shards: Shard{Total: 2, Successful: 1, Failed: 0}, Created: true},
		Success{Acknowledged: true},
	}
	// check parsing result
	checkParsingResult(t, input, parser, expected)
}

// calculated actual parsing result and check it against expected result
func checkParsingResult(t *testing.T, input []string, parser Parser, expected []interface{}) {
	// actual result from parsing input
	var actual []interface{}
	for _, in := range input {
		bytes := []byte(in)
		result, _ := parser.Parse(bytes)
		actual = append(actual, result)
	}
	equalsInterface(t, actual, expected)
}
