package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap19 runs example queries from chapter 19 of Elasticsearch the Definitive Guide
// It's an introduction to identifying words (i.e. tokenizer) using standard analyzer regardless of the document language.
func chap19() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// using 'whitespace' tokenizer
	c.Analyze("").AddParam(e.Tokenizer, "whitespace").Get("You're the 1st runner home!")
	// using 'letter' tokenizer
	c.Analyze("").AddParam(e.Tokenizer, "letter").Get("You're the 1st runner home!")
	// using 'letter' tokenizer
	c.Analyze("").AddParam(e.Tokenizer, "standard").Get("You're the 1st runner home!")
	// using 'uax_url_email' tokenizer
	c.Analyze("").AddParam(e.Tokenizer, "uax_url_email").Get("You should send him a message on john@doe.dz")

	// set character filter to clean HTML text
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_html_analyzer", e.Dict{e.Tokenizer: "standard", e.CharFilters: []string{"html_strip"}})).Put()
	// check the html text is well preprocessed
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("my_html_analyzer").Pretty().Get(`<p>Some d&eacute;j&agrave; vu <a href="http://somedomain.com">website</a>`)

	// replace the different alternative for apostrophe with a single one
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("char_filter").Add2("quotes", e.Dict{"type": "mapping", "mappings": []string{"\\u0091=>\\u0027", "\\u0092=>\\u0027", "\\u2018=>\\u0027", "\\u2019=>\\u0027", "\\u201B=>\\u0027"}})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("quotes_analyzer", e.Dict{e.Tokenizer: "standard", e.CharFilters: []string{"quotes"}})).Put()
	// test the analyzer after creating it
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("quotes_analyzer").Get("You‘re my ’favorite‘ M'Coy")
}
