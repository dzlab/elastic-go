package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap24 runs example queries from chapter 24 of Elasticsearch the Definitive Guide
// It's about fuzzy matching at query-time to hadling typoes and mispelings, and phonetic token filters at index time for sounds-like matching.
func chap24() {
	c := &Elasticsearch{Addr: "localhost:9200"}
	c.Index("my_index").Delete()
	// index some documents
	c.Bulk("my_index", "my_type").AddOperation(e.NewOperation(1).Add("text", "Surprise me!")).AddOperation(e.NewOperation(2).Add("text", "That was surprising.")).AddOperation(e.NewOperation(3).Add("text", "I wasn't surprised.")).Post()
	// fuzzy query
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewFuzzyQuery().Add("text", "surprize"))).Get()
	// set the fuziness to limit number of matching documents
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewFuzzyQuery().Add("text", "surprize").Add(e.Fuzziness, 1))).Get()

	// 'match' query with fuzziness: query text will be analyzed and each term will be fuzzified
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("text").Add("query", "SURPRIZE ME!").Add(e.Fuzziness, 1).Add(e.Operator, "and")))).Get()
	// 'multi_match' query with fuzziness
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "SURPRIZE ME!").Add(e.Fuzziness, "AUTO").Add("fields", []string{"text", "title"}))).Get()

	// using phonetic plugin
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("dbl_metaphone", e.Dict{e.Type: "phonetic", e.Encoder: "double_metaphone"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("dbl_metaphone", e.Dict{e.Tokenizer: "standard", "filter": "dbl_metaphone"})).Put()
	// Test the analyzer
	c.Analyze("my_index").Analyzer("dbl_metaphone").Get("Smith Smythe")
	// use the phonetic analyzer in a document mapping definition
	c.Mapping("my_index", "my_type").AddField("name", e.Dict{e.Type: "string", "fields": e.Dict{"phonetic": e.Dict{"type": "string", "analyzer": "dbl_metaphone"}}}).Put()
	// insert some documents
	c.Insert("my_index", "my_type").Document(1, e.Dict{"name": "John Smith"}).Put()
	c.Insert("my_index", "my_type").Document(2, e.Dict{"name": "Jonnie Smythe"}).Put()
	t.Sleep(1 * t.Second)
	// now search: both document should be returned
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("name.phonetic").Add("query", "Jahnnie Smeeth").Add("operator", "and")))).Get()
}
