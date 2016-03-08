package main

import (
	"github.com/dzlab/elastic-go"
)

func main() {
	client := &elastic.Elasticsearch{Addr: "localhost:9200"}
	client.Search("", "").Add("from", "30").Add("size", "10").Get()
	client.Search("", "").AddQuery(elastic.NewQuery("query").AddQuery(elastic.NewQuery("match").Add("tweet", "elasticsearch"))).Get()
	client.Search("index_2014*", "type1,type2").Get()
	client.Index("gb").Delete()
	client.Index("gb").Put("{mappings: {tweet: {properties: {tweet:{type: \"string\", analyzer: \"english\"}, date: {type: \"date\"}, name: {type: \"string\"}, user_id: {type: \"long\"}}}}}")
	client.Mapping("gb", "tweet").Put("{properties: {tag: {type: \"string\", index: \"not_analyzed\"}}}")
	client.Mapping("gb", "tweet").Get()
	client.Analyze("gb").Get("tweet")
	client.Analyze("gb").Get("tag")
}
