package main

import (
	e "github.com/dzlab/elastic-go"
)

func main() {
	client := &e.Elasticsearch{Addr: "localhost:9200"}
	client.Search("", "").Add("from", "30").Add("size", "10").Get()
	client.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").Add("tweet", "elasticsearch"))).Get()
	client.Search("index_2014*", "type1,type2").Get()
	client.Index("gb").Delete()
	//"{mappings: {tweet: {properties: {tweet:{type: \"string\", analyzer: \"english\"}, date: {type: \"date\"}, name: {type: \"string\"}, user_id: {type: \"long\"}}}}}"
	client.Index("gb").Put()
	client.Mapping("gb", "tweet").Put("{properties: {tag: {type: \"string\", index: \"not_analyzed\"}}}")
	client.Mapping("gb", "tweet").Get()
	client.Analyze("gb").Get("tweet")
	client.Analyze("gb").Get("tag")
	client.Validate("gb", "tweet", true).AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("tweet").Add("match", "really powerful"))).Get()
	// sort by date
	client.Search("", "").AddQuery(
		e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("filter").AddQuery(e.NewQuery("term").Add("user_id", "1")))),
	).AddQuery(
		e.NewQuery("sort").AddQuery(e.NewQuery("date").Add("order", "desc")),
	).Get()
}
