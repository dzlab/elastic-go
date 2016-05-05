package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap22 runs example queries from chapter 22 of Elasticsearch the Definitive Guide
// It's about using stopwords without compromising perforance or precision
func chap22() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// to use custom words, create an anlyzer and pass an array of stopwords or language specific list (e.g. _english_)
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_analyzer", e.Dict{"type": "standard", e.Stopwords: []string{"and", "the"}})).Put()
	// test the new analyzer with some text
	t.Sleep(500 * t.Millisecond)
	c.Analyze("my_index").Analyzer("my_analyzer").Get("The quick and the dead")

	// specify a file containing stopwords
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_analyzer", e.Dict{"type": "standard", e.StopwordsPath: "stopwords/english.txt"})).Put()
	// test the new analyzer with some text
	t.Sleep(500 * t.Millisecond)
	c.Analyze("my_index").Analyzer("my_analyzer").Get("The quick and the dead")

	// disabling stopwords usage
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_english", e.Dict{"type": "english", "stopwords": "_none_"})).Put()
	t.Sleep(500 * t.Millisecond)
	c.Analyze("my_index").Analyzer("my_english").Get("The quick and the dead")

	// create a custom spanish analyzer
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("spanish_stop", e.Dict{"type": "stop", e.Stopwords: []string{"si", "esta", "el", "la"}}).Add2("light_spanish", e.Dict{"type": e.Stemmer, e.Language: "light_spanish"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_spanish", e.Dict{e.Tokenizer: "spanish", "filter": []string{"lowercase", "asciifolding", "spanish_stop", "light_spanish"}})).Put()

	// reduce index size with index_options to only information about which document contanins which terms and frequency of terms in those documents (i.e. freqs option)
	c.Index("my_index").Mappings("my_type", e.NewMapping().AddField("title", e.Dict{e.Type: "string"}).AddField("content", e.Dict{e.Type: "string", e.IndexOptions: "freqs"})).Put()

	// example of using common_grams token filter
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("index_filter", e.Dict{e.Type: e.CommonGrams, e.CommonWords: "_english_"}).Add2("search_filter", e.Dict{e.Type: e.CommonGrams, e.CommonWords: "_english", e.QueryMode: true})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("index_grams", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "index_filter"}}).Add2("search_grams", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "search_filter"}})).Put()
	// then create a filed that used index_grams at index time
	c.Mapping("my_index", "my_type").AddField("text", e.Dict{e.Type: "string", e.IndexAnalyzer: "index_grams", e.SearchAnalyzer: "standard"}).Put()
	// use the index_grams to analyze some text (you will see bigrams), and search_grams
	c.Analyze("my_index").Analyzer("index_grams").Get("The quick and brwon fox")
	c.Analyze("my_index").Analyzer("search_grams").Get("The quick and brwon fox")
	// the index contains unigrams so it can be queried as usual
	c.Search("my_index", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("text").Add("query", "the quick and brown fox").Add(e.CutOffFrequency, 0.01)))).Get()
	// bigram phrase queries
	c.Search("my_index", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("text").Add("query", "The quick and brown fox").Add("analyzer", "search_grams")))).Get()
	// Two word phrases become faster as it turns out to be a single term search
	c.Search("my_index", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("text").Add("query", "The quick").Add("analyzer", "search_grams")))).Get()
}
