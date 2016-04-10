package elastic

import (
//"encoding/json"
)

/*
 * Elasticsearch failure representation
 * e.g.:{"error":{"root_cause":[{"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"}],"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"},"status":503}
 * e.g.:{"error":{"root_cause":[{"type":"index_already_exists_exception","reason":"already exists","index":"my_index"}],"type":"index_already_exists_exception","reason":"already exists","index":"my_index"},"status":400}
 */
type Failure struct {
	err Error `json:"error"`
}
type Error struct {
	rootCause []Dict `json:"root_cause"`
	kind      string `json:"type"`
	reason    string `json:"reason"`
	status    int    `json:"status"`
}

/*
 * Elasticsearch success representation
 * e.g.: {"acknowledged":true}
 */
type Success struct {
	Acknowledged bool `json:"acknowledged"`
}

/*
 * Elasticsearch unvalid representation
 * e.g.: {"valid":false,"_shards":{"total":1,"successful":1,"failed":0},"explanations":[{"index":"gb","valid":false,"error":"org.elasticsearch.index.query.QueryParsingException: No query registered for [tweet]"}]}
 */
type Unvalid struct {
	Valid       bool   `json:"valid"`
	Shards      Dict   `json:"_shards"`
	Explanation []Dict `json:"explanations"`
}

/////////////////////////////////// Search Query

/*
 * Elasticsearch shard response representation
 */
type Shard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

/*
 * Elasticsearch search hits representation
 */
type Hits struct {
	Total    int         `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits     []Dict      `json:"hits"`
}

/*
 * Elastisearch search result representation
 * e.g. {"took":1,"timed_out":false,"_shards":{"total":5,"successful":5,"failed":0},"hits":{"total":0,"max_score":null,"hits":[]}}
 */
type SearchResult struct {
	Took     int   `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Shards   Shard `json:"_shards"`
	Hits     Hits  `json:"hits"`
}

/////////////////////////////////// Analyze Query

/*
 * Elasticsearch analyze result representation
 * e.g. {"tokens":[{"token":"quick","start_offset":0,"end_offset":5,"type":"<ALPHANUM>","position":0},{"token":"brown","start_offset":6,"end_offset":11,"type":"<ALPHANUM>","position":1},{"token":"fox","start_offset":12,"end_offset":15,"type":"<ALPHANUM>","position":2}]}
 */

/////////////////////////////////// Insert Query

/*
 * Elasticsearch insert result representation
 * e.g. {"_index":"my_index","_type":"groups","_id":"1","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"created":true}
 */

/////////////////////////////////// Bulk Query
/*
 * Elasticsearch bulk result representation
 * e.g. {"took":118,"errors":false,"items":[{"index":{"_index":"my_index","_type":"my_type","_id":"1","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"status":201}},{"index":{"_index":"my_index","_type":"my_type","_id":"2","_version":1,"_shards":{"total":2,"successful":1,"failed":0},"status":201}}]}
 */
