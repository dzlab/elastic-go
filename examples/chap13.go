package main

import (
	e "github.com/dzlab/elastic-go"
	"time"
)

/*
 * Examples of queries based on Elasticsearch Definitive Guide, chapter 13
 * Full text search examples
 */
func chap13() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// create index, delete if any
	c.Index("my_index").Delete()
	// use one primary shard
	c.Index("my_index").AddSetting(e.SHARDS_NB, 1).Put()

	// insert some test data
	c.Bulk("my_index", "my_type").AddOperation(e.NewOperation(1).Add("title", "The quick brown fox")).AddOperation(e.NewOperation(2).Add("title", "The quick brown fox jumps over the lazy dog")).AddOperation(e.NewOperation(3).Add("title", "The quick brown fox jumps over the quick dog")).AddOperation(e.NewOperation(4).Add("title", "Brown fox brown dog")).Post()

	// after a bulk insert, we have to wait for the inserted documents to be available
	time.Sleep(1 * time.Second)

	// single word query
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").Add("title", "QUICK!"))).Get()

	// multi-word queries: brown OR dog
	// i.e. any document whose title match at least one field is returned
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").Add("title", "BROWN DOG!"))).Get()

	// improving precision: brown AND dog
	// i.e. exclude documents that contains only one term
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").AddQuery(e.NewQuery("title").Add("query", "BROWN DOG!").Add("operator", "and")))).Get()

	// controlling precision:
	// i.e. include documents that contains at least one and/or 75% of query terms
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").AddQuery(e.NewQuery("title").Add("query", "quick brown dog").Add("minimum_should_match", "75%")))).Get()

	// combining queries: `bool` queries
	// i.e. just like `bool` filters but score documents by relevance
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddMust(e.NewQuery("match").Add("title", "quick")).AddMustNot(e.NewQuery("match").Add("title", "lazy")).AddShould(e.NewQuery("match").Add("title", "brown")).AddShould(e.NewQuery("match").Add("title", "dog")))).Get()

	// `should` clauses are not supposed to match, but if there is no `must` than at least one `should` query have to match
	// in the following request, at least 2 terms have to be in the document be returned
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").Add("title", "brown")).AddShould(e.NewQuery("match").Add("title", "fox")).AddShould(e.NewQuery("match").Add("title", "dog")).Add("minimum_should_match", 2))).Get()

	// look for documents with `full` and `text` and `search`, give those containing `Elasticsearch` and `Lucene` higher score
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddMust(e.NewQuery("match").AddQuery(e.NewQuery("content").Add("query", "full text search").Add("operator", "and"))).AddShould(e.NewQuery("match").Add("content", "Elasticsearch")).AddShould(e.NewQuery("match").Add("content", "Lucene")))).Get()
	// boost relevance of documents containing 'Elasticsearch' over `Lucene`
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddMust(e.NewQuery("match").AddQuery(e.NewQuery("content").Add("query", "full text search").Add("operator", "and"))).AddShould(e.NewQuery("match").AddQuery(e.NewQuery("content").Add("query", "Elasticsearch").Add("boost", 3))).AddShould(e.NewQuery("match").AddQuery(e.NewQuery("content").Add("query", "Lucene").Add("boost", 2))))).Get()

	// controling text analysis
	// add a field to document
	c.Mapping("my_index", "my_type").AddProperty("english_title", "type", "string").AddProperty("english_title", "analyzer", "english").Put()

	// compare how values in `english_title` and `title` are analyzed at index time
	c.Analyze("my_index").Field("my_type.title").Get("Foxes")
	c.Analyze("my_index").Field("my_type.english_title").Get("Foxes")

	// check how `match` query is analyzed (hint `english_title` uses english analyzer)
	c.Validate("my_index", "my_type", true).AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").Add("title", "Foxes")).AddShould(e.NewQuery("match").Add("english_title", "Foxes")))).Get()
}
