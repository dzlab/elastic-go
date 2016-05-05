package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap21 runs example queries from chapter 21 of Elasticsearch the Definitive Guide
// It's about stemming, it.e. reducing words to their root form
func chap21() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// delete index if exists
	c.Index("my_index").Delete()
	// custom analyzer based on the English Algorithmic stemmer configuration
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("english_stop", e.Dict{e.Type: "stop", e.Stopwords: "_english_"}).Add2("light_english_stemmer", e.Dict{e.Type: e.Stemmer, e.Language: "english_light"}).Add2("english_possessive_stemmer", e.Dict{e.Type: e.Stemmer, e.Language: "possessive_english"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("english", e.Dict{e.Tokenizer: "standard", e.FILTER: []string{"english_possessive_stemmer", "lowercase", "english_stop", "light_english_stemmer", "asciifolding"}})).Put()
	// Dictionary stemmer based on hunspell (install hunspell p361)
	// Create a hunspell-based filter
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("en_US", e.Dict{e.Type: "hunspell", e.Language: "en_US"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("en_US", e.Dict{e.TOKENIZER: "standard", "filter": []string{"lowercase", "en_US"}})).Put()
	// analyze some text use standard analyzer and the en_US custom analyzer
	c.Analyze("my_index").Analyzer("en_US").Get("reorganizes")
	c.Analyze("").Analyzer("english").Get("reorganizes")

	// customize stemming by preventing it for some words
	// Instead of an array of words for 'keywords' we could also specify a file with 'keywords_path'
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("no_stem", e.Dict{e.Type: "keyword_marker", "keywords": []string{"skies"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_english", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "no_stem", "porter_stem"}})).Put()
	// test the custom analyzer to check how 'ski' words are stemmed
	c.Analyze("my_index").Analyzer("my_english").Get("sky skies skiing skis")

	// we can also specify how a word to be stemmed with 'stemmer_override'
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("custom_stem", e.Dict{e.TYPE: "stemmer_override", "rules": []string{"skies=>sky", "mice=>mouse", "feet=>foot"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_english", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "custom_stem", "porter_stem"}})).Put()
	// now we can test this custom analyzer with some text to verity the correct stemming
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("my_english").Get("The mice came down from the skies and ran over my feet")

	// stemming in situ: indexing the original word and its stemmed form
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("unique_stem", e.Dict{"type": "unique", "only_on_same_position": true})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("in_situ", e.Dict{e.Tokenizer: "standard", "filter": []string{"lowercase", "keyword_repeat", "porter_stem", "unique_stem"}})).Put()
	// it is recommended not to use stemming in situ, p. 370
}
