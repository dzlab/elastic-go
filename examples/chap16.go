package main

import (
	e "github.com/dzlab/elastic-go"
	"time"
)

func chap16() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("address", e.NewMapping().AddProperty("postcode", "type", "string").AddProperty("postcode", "index", "not_analyzed")).Put()

	c.Insert("my_index", "address").Document(1, e.Dict{"postcode": "W1V 3DG"}).Put()
	c.Insert("my_index", "address").Document(2, e.Dict{"postcode": "W2F 8HW"}).Put()
	c.Insert("my_index", "address").Document(3, e.Dict{"postcode": "W1F 7HW"}).Put()
	c.Insert("my_index", "address").Document(4, e.Dict{"postcode": "WC1N 1LZ"}).Put()
	c.Insert("my_index", "address").Document(5, e.Dict{"postcode": "SW5 0BE"}).Put()

	// wait for data to be available
	time.Sleep(1 * time.Second)

	// prefix query
	c.Search("my_index", "address").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.Prefix).Add("postcode", "W1"))).Get()

	// wildcard query: ? match any character, * matches zero or more
	c.Search("my_index", "address").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.Wildcard).Add("postcode", "W?F*HW"))).Get()
	// regular expression query
	c.Search("my_index", "address").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.RegExp).Add("postcode", "W[0-9].+"))).Get()

	// phrase query with prefix can be used for instant search (i.e. returning results to users as they are typing)
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.MatchPhrasePrefix).Add("brad", "johnie walker bl"))).Get()

	// control the order of terms in phrase query
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.MatchPhrasePrefix).AddQuery(e.NewQuery("brand").Add("query", "johnie walker bl").Add(e.SLOP, 10)))).Get()
	// control how many terms the prefix can be expanded to
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery(e.MatchPhrasePrefix).AddQuery(e.NewQuery("brand").Add("query", "johnie walker bl").Add(e.MaxExpansions, 50)))).Get()

	// index-time search as you type
	c.Index("my_index").Delete()
	c.Index("my_index").SetShardsNb(1).AddAnalyzer(e.NewAnalyzer("filter").Add2("autocomplete_filter", e.Dict{"type": "edge_ngram", "min_gram": 1, "max_gram": 20})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("autocomplete", e.Dict{"type": "custom", "tokenizer": "standard", "filter": []string{"lowercase", "autocomplete_filter"}})).Put()
	// we can analyze the behavior of this analyzer
	time.Sleep(1 * time.Second)
	c.Analyze("my_index").Analyzer("autocomplete").Get("quick brown")
	// in order to use the analyzer, we need to apply it to a field
	c.Mapping("my_index", "my_type").AddProperty("name", "type", "string").AddProperty("name", "analyzer", "autocomplete").Put()

	c.Bulk("my_index", "my_type").AddOperation(e.NewOperation(1).Add("name", "Brown foxes")).AddOperation(e.NewOperation(2).Add("name", "Yellow furballs")).Post()

	time.Sleep(1 * time.Second)
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("name", "brown fo"))).Get()
	// the validation api shine some lights to understand the search result
	c.Validate("my_index", "my_type", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("name", "brown fo"))).Get()
	// we can overide at query time the autocomplete analyzer which has been used at index and query time
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("name").Add("query", "brown fo").Add("analyzer", "standard")))).Get()
	// alternatively we can sepcify separate anlyzers for index and search time by a mapping
	c.Mapping("my_index", "my_type").AddField("name", e.Dict{"type": "string", e.IndexAnalyzer: "autocomplete", e.SearchAnalyzer: "standard"}).Put()
	// repeat the previous validate query
	c.Validate("my_index", "my_type", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("name", "brown fo"))).Get()

	// use keyword tokenizer (that do nothing) as postcode needs to be analyzed
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("postcode_filter", e.Dict{"type": "edge_ngram", "min_gram": 1, "max_gram": 8})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("postcode_index", e.Dict{"tokenizer": "keyword", "filter": []string{"postcode_filter"}}).Add1("postcode_search", "tokenizer", "keyword")).Put()

	// for searching compound words in languages like German, trigram is a good starting point
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("trigrams_filter", e.Dict{"type": "ngram", "min_gram": 3, "max_gram": 3})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("trigrams", e.Dict{"type": "custom", "tokenizer": "standard", "filter": []string{"lowercase", "trigrams_filter"}})).Mappings("my_type", e.NewMapping().AddField("text", e.Dict{"type": "string", "analyzer": "trigrams"})).Put()
	// test the trigrams analyzer
	time.Sleep(1 * time.Second)
	c.Analyze("my_index").Analyzer("trigrams").Get("Weibkopfseeadler")
	// post some data
	c.Bulk("my_index", "my_type").AddOperation(e.NewOperation(1).Add("text", "Aussprachewörterbuch")).AddOperation(e.NewOperation(2).Add("text", "Militärgeschichte")).AddOperation(e.NewOperation(3).Add("text", "Wiebkopfseeadler")).AddOperation(e.NewOperation(4).Add("text", "Weltgesundheitsorganisation")).AddOperation(e.NewOperation(1).Add("text", "Rindfleischetikettierungsüberwachungsaufgabenübertragungsgesetz")).Post()
	// search
	time.Sleep(1 * time.Second)
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("text", "Adler"))).Get()
	// use minimum_should_match to remove spurius results
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("text").Add("query", "Adler").Add(e.MinimumShouldMatch, "80%")))).Get()
}
