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
		`{"error":{"root_cause":[{"type":"mapper_parsing_exception","reason":"failed to parse [title]"}],"type":"mapper_parsing_exception","reason":"failed to parse [title]","caused_by":{"type":"json_parse_exception","reason":"Unexpected end-of-input in VALUE_STRING\n at [Source: org.elasticsearch.common.io.stream.InputStreamStreamInput@1281d55b; line: 1, column: 35]"}},"status":400}`,
	}
	// expected results
	expected := []interface{}{
		InsertResult{Index: "blogposts", Doctype: "post", ID: "1", Version: 2, Shards: Shard{Total: 2, Successful: 1, Failed: 0}, Created: false},
		InsertResult{Index: "my_index", Doctype: "groups", ID: "1", Version: 1, Shards: Shard{Total: 2, Successful: 1, Failed: 0}, Created: true},
		Success{Acknowledged: true},
		Failure{Err: Error{RootCause: []Dict{Dict{"type": "mapper_parsing_exception", "reason": "failed to parse [title]"}}, Type: "mapper_parsing_exception", Reason: "failed to parse [title]", CausedBy: Dict{"type": "json_parse_exception", "reason": "Unexpected end-of-input in VALUE_STRING\n at [Source: org.elasticsearch.common.io.stream.InputStreamStreamInput@1281d55b; line: 1, column: 35]"}}, Status: 400},
	}
	// check parsing result
	checkParsingResult(t, input, parser, expected)
}

// TestIndexResultParser tests for IndexResultParser
func TestIndexResultParser(t *testing.T) {
	parser := &IndexResultParser{}
	// input strings to parse
	input := []string{
		`{"error":{"root_cause":[{"type":"index_already_exists_exception","reason":"already exists","index":"my_index"}],"type":"index_already_exists_exception","reason":"already exists","index":"my_index"},"status":400}`,
		`{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index","resource.type":"index_or_alias","resource.id":"my_index","index":"my_index"}],"type":"index_not_found_exception","reason":"no such index","resource.type":"index_or_alias","resource.id":"my_index","index":"my_index"},"status":404}`,
	}
	// expected results
	expected := []interface{}{
		Failure{Err: Error{RootCause: []Dict{Dict{"type": "index_already_exists_exception", "reason": "already exists", "index": "my_index"}}, Type: "index_already_exists_exception", Reason: "already exists", Index: "my_index"}, Status: 400},
		Failure{Err: Error{RootCause: []Dict{Dict{"type": "index_not_found_exception", "reason": "no such index", "resource.type": "index_or_alias", "resource.id": "my_index", "index": "my_index"}}, Type: "index_not_found_exception", Reason: "no such index", ResourceType: "index_or_alias", ResourceId: "my_index", Index: "my_index"}, Status: 404},
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
