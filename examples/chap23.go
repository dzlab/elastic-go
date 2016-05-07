package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap23 runs example queries from chapter 23 of Elasticsearch the Definitive Guide.
// It's about how to handle synonyms
func chap23() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// create synonym tokens as follows
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("my_synonym_filter", e.Dict{e.Type: "synonym", e.Synonyms: []string{"british,english", "queen,monarch"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_synonyms", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "my_synonym_filter"}})).Put()
	// test the analyzer with some text
	c.Analyze("my_index").Analyzer("my_synonyms").Get("Elizabeth is the English queen.")

	// Some bizare stuff with multi-word synonyms and phrase queries
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("my_synonym_filter", e.Dict{e.Type: "synonym", e.Synonyms: []string{"usa,united states,u s a,united states of america"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_synonyms", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "my_synonym_filter"}})).Put()
	c.Analyze("my_index").Analyzer("my_synonyms").Get("The United States is wealthy")
	// even more bizare stuff when using synonyms at query time
	c.Validate("my_index", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("text").Add("query", "usa is wealthy").Add("analyzer", "my_synonyms")))).Get()
	// to avoid this mess, use 'simple contraction' (which simply '=>' rule) for phrase queries
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("my_synonym_filter", e.Dict{e.Type: "synonym", e.Synonyms: []string{"united states,u s a,united states of america=>usa"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_synonyms", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "my_synonym_filter"}})).Put()
	// now test with same previous text
	c.Analyze("my_index").Analyzer("my_synonyms").Get("The United States is wealthy")
	c.Validate("my_index", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("text").Add("query", "usa is wealthy").Add("analyzer", "my_synonyms")))).Get()

	// use char mapping filter to convert emoticons to their meaning and void to lose them as the standard tokenizer filter will remove them.
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("char_filter").Add2("emoticons", e.Dict{e.Type: "mapping", e.MAPPINGS: []string{":)=>emoticon_happy", ":(=>emoticon_sad"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_emoticons", e.Dict{e.CharFilter: "emoticons", e.Tokenizer: "standard", "filter": []string{"lowercase"}})).Get()
	// test the symbol synonym
	c.Analyze("my_index").Analyzer("my_emoticons").Get("I am :) not :(")
}
