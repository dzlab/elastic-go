package elastic

import (
//"encoding/json"
)

// Failure is a structure representing the Elasticsearch failure response
// e.g.:{"error":{"root_cause":[{"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"}],"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"},"status":503}
// e.g.:{"error":{"root_cause":[{"type":"index_already_exists_exception","reason":"already exists","index":"my_index"}],"type":"index_already_exists_exception","reason":"already exists","index":"my_index"},"status":400}
type Failure struct {
	err Error `json:"error"`
}

// Error is a structure representing the Elasticsearch error response
type Error struct {
	rootCause []Dict `json:"root_cause"`
	kind      string `json:"type"`
	reason    string `json:"reason"`
	status    int    `json:"status"`
}

// Success is a structure representing an Elasticsearch success response
// e.g.: {"acknowledged":true}
type Success struct {
	Acknowledged bool `json:"acknowledged"`
}

// Unvalid is a structure representing an Elasticsearch unvalid response
// e.g.: {"valid":false,"_shards":{"total":1,"successful":1,"failed":0},"explanations":[{"index":"gb","valid":false,"error":"org.elasticsearch.index.query.QueryParsingException: No query registered for [tweet]"}]}
type Unvalid struct {
	Valid       bool   `json:"valid"`
	Shards      Dict   `json:"_shards"`
	Explanation []Dict `json:"explanations"`
}

/////////////////////////////////// Search Query

// Shard is a structure representing the Elasticsearch shard part of Search query response
type Shard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// Hits is a structure representing the Elasticsearch hits part of Search query response
type Hits struct {
	Total    int         `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits     []Dict      `json:"hits"`
}

/*
 * Elasticsearch explain result
 * e.g. {"valid":true,"_shards":{"total":1,"successful":1,"failed":0},"explanations":[{"index":"my_index","valid":true,"explanation":"+((name:b name:br name:bro name:brow name:brown) (name:f name:fo)) #ConstantScore(+ConstantScore(_type:my_type))"}]}
 */

// SearchResult is a structure representing the Elastisearch search result
// e.g. {"took":1,"timed_out":false,"_shards":{"total":5,"successful":5,"failed":0},"hits":{"total":0,"max_score":null,"hits":[]}}
// e.g. {"took":3,"timed_out":false,"_shards":{"total":1,"successful":1,"failed":0},"hits":{"total":1,"max_score":0.50741017,"hits":[{"_index":"my_index","_type":"my_type","_id":"1","_score":0.50741017,"_source":{"name":"Brown foxes"}}]}}
type SearchResult struct {
	Took     int   `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Shards   Shard `json:"_shards"`
	Hits     Hits  `json:"hits"`
}

/////////////////////////////////// Analyze Query

// AnalyzeResult is a structure representing the Elasticsearch analyze query result
// e.g. {"tokens":[{"token":"quick","start_offset":0,"end_offset":5,"type":"<ALPHANUM>","position":0},{"token":"brown","start_offset":6,"end_offset":11,"type":"<ALPHANUM>","position":1},{"token":"fox","start_offset":12,"end_offset":15,"type":"<ALPHANUM>","position":2}]}
type AnalyzeResult struct {
	tokens []AnalyzeToken `json:"tokens"`
}

// AnalyzeToken is a structure representing part of the Elasticsearch analyze query response
type AnalyzeToken struct {
	token       string `json:"token"`
	startOffset int    `json:"start_offset"`
	endOffset   int    `json:"end_offset"`
	tokenType   string `json:"type"`
	position    int    `json:"position"`
}

/////////////////////////////////// Insert Query

// InsertResult is a strucuture representing the Elasticsearch insert query result
// e.g. {"_index":"my_index","_type":"groups","_id":"1","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"created":true}
type InsertResult struct {
	index   string `json:"_index"`
	doctype string `json:"_type"`
	id      string `json:"_id"`
	version int    `json:"version"`
	shards  Shard  `json:"_shards"`
	created bool   `json:"created"`
	status  int    `json:"status"`
}

/////////////////////////////////// Bulk Query

// BulkResult is a structure representing the Elasticsearch bulk query result
// e.g. {"took":118,"errors":false,"items":[{"index":{"_index":"my_index","_type":"my_type","_id":"1","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"status":201}},{"index":{"_index":"my_index","_type":"my_type","_id":"2","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"status":201}}]}
type BulkResult struct {
	took   int            `josn:"took"`
	errors bool           `json:"errors"`
	items  []InsertResult `json:"items"`
}
