package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap20 runs example queries from chapter 20 of Elasticsearch the Definitive Guide
// It's about normalizing tokens, e.g. lowercase
func chap20() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// using the lowercase filter
	c.Analyze("").AddParam(e.Tokenizer, "standard").AddParam("filters", "lowercase").Get("The QUICK Brown FOX!")

	// use token filters as part of the analysis process
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_lowercase", e.Dict{e.Tokenizer: "standard", e.Filter: []string{"lowercase"}})).Put()
	// test the analyzer
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("my_lowercase").Get("The QUICK Brown FOX!")

	// Diacritics (e.g. ‘, ^ and ¨) can be cleaned with asciifolding filter
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("folding", e.Dict{e.Tokenizer: "standard", e.Filter: []string{"lowercase", "asciifolding"}})).Put()
	// test the analyzer
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("folding").Get("My œsophagus caused a débâcle")

	// removing diacritics may cause the loose of meaning, to avoid this index the text twice: once in the original form and once with diacritics removed
	c.Mapping("my_index", "my_type").AddField("title", e.Dict{"type": "string", "analyzer": "standard", "fields": e.Dict{"folded": e.Dict{"type": "string", "analyzer": "folding"}}}).Put()
	// test mapping
	c.Analyze("my_index").Field("title").Get("Esta està loca")
	c.Analyze("my_index").Field("title.folded").Get("Esta està loca")
	// insert document for further testing
	c.Insert("my_index", "my_type").Document(1, e.Dict{"title": "Esta loca!"}).Put()
	c.Insert("my_index", "my_type").Document(2, e.Dict{"title": "Està loca!"}).Put()
	t.Sleep(1 * t.Second)
	c.Search("my_index", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("type", "most_fields").Add("query", "està loca").AddMultiple("fields", "title", "title.folded"))).Get()
	// Explain the query for better understanding
	t.Sleep(1 * t.Second)
	c.Validate("my_index", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("type", "most_fields").Add("query", "està loca").AddMultiple("fields", "title", "title.folded"))).Get()

	// use icu_normalizer token filter to ensure that all tokens are in the same form
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("nfkc_normalizer", e.Dict{"type": "icu_normalizer", "name": "nfkc"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_normalizer", e.Dict{e.Tokenizer: "icu_tokenizer", "filter": []string{"nfkc_normalizer"}})).Put()

	// nfkc_cf is the equivalent of lowercase token filter but suitable to all languages
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_lowercaser", e.Dict{e.Tokenizer: "icu_tokenizer", "filter": []string{"icu_normalizer"}})).Put()
	t.Sleep(1 * t.Second)
	c.Analyze("").Analyzer("standard").Get("Weißkopfseeadler WEISSKOPFSEEADLER")

	// Unicode character folding
	c.Analyze("my_index").Analyzer("my_lowercaser").Get("Weißkopfseeadler WEISSKOPFSEEADLER")

	// the icu_folding token filter applies automatically Unicode normalization and case folding nfkc_cf, no need for icu_normalizer
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_folder", e.Dict{e.Tokenizer: "icu_tokenizer", e.Filter: []string{"icu_folding"}})).Put()
	// test the analyzer, see how arabic numeral are folded to latin equivalent
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Analyzer("my_folder").Get("")

	// we can prevent some characters from been folded
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("filter").Add2("swedish_folding", e.Dict{"type": "icu_folding", "unicodeSetFilter": "[^äöÄÖ]"})).AddAnalyzer(e.NewAnalyzer("filter").Add2("swedish_analyzer", e.Dict{e.Tokenizer: "icu_tokenizer", "filter": []string{"swedish_folding", "lowercase"}})).Put()

	// sorting and collations
	// case insensitive sorting
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("user", e.NewMapping().AddField("name", e.Dict{"type": "string", "fields": e.Dict{"raw": e.Dict{"type": "string", "index": "not_analyzed"}}})).Put()
	c.Insert("my_index", "user").Document(1, e.Dict{"name": "Boffey"}).Put()
	c.Insert("my_index", "user").Document(2, e.Dict{"name": "BROWN"}).Put()
	c.Insert("my_index", "user").Document(3, e.Dict{"name": "bailey"}).Put()
	t.Sleep(1 * t.Second)
	// sort the names lexicographically
	c.Search("my_index", "user").AddParam("sort", "name.raw").Get()

	// to sort the names in the more natuarl alphabetically use a subfield that is lowercase of the original 'name' field
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("case_insensitive_sort", e.Dict{e.Tokenizer: "keyword", "filter": []string{"lowercase"}})).Put()
	c.Mapping("my_index", "user").AddField("name", e.Dict{"type": "string", "fields": e.Dict{"lower_case_sort": e.Dict{"type": "string", "analyzer": "case_insensitive_sort"}}}).Put()
	c.Insert("my_index", "user").Document(1, e.Dict{"name": "Boffey"}).Put()
	c.Insert("my_index", "user").Document(2, e.Dict{"name": "BROWN"}).Put()
	c.Insert("my_index", "user").Document(3, e.Dict{"name": "bailey"}).Put()
	t.Sleep(1 * t.Second)
	c.Search("my_index", "user").AddParam("sort", "name.lower_case_sort").Get()

	// sorting for many languages using icu_collation token filter default sorting DUCET
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("ducet_sort", e.Dict{e.Tokenizer: "keyword", "filter": []string{"icu_collation"}})).Put()
	c.Mapping("my_index", "user").AddField("name", e.Dict{"type": "string", "fields": e.Dict{"sort": e.Dict{"type": "string", "analyzer": "ducet_sort"}}}).Put()
	// as we didn't specify a language, it's defaults to using DUCET collation
	c.Bulk("my_index", "user").AddOperation(e.NewOperation(1).Add("name", "Boffey")).AddOperation(e.NewOperation(2).Add("name", "BROWN")).AddOperation(e.NewOperation(3).Add("name", "bailey")).AddOperation(e.NewOperation(4).Add("name", "Böhm")).Post()
	t.Sleep(1 * t.Second)
	c.Search("my_index", "user").AddParam("sort", "name.sort").Get()

	// the icu_collation filter can be configured with a language
	// e.g. setup German phonebook sort order
	c.Index("my_index").Delete()
	c.Index("my_index").SetShardsNb(1).AddAnalyzer(e.NewAnalyzer("filter").Add2("german_phonebook", e.Dict{"type": "icu_collation", "language": "de", "country": "DE", "variant": "@collation=phonebook"})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("german_phonebook", e.Dict{e.Tokenizer: "keyword", "filter": []string{"german_phonebook"}})).Mappings("user", e.NewMapping().AddField("name", e.Dict{"type": "string", "fields": e.Dict{"sort": e.Dict{"type": "string", "analyzer": "german_phonebook"}}})).Put()
	c.Bulk("my_index", "user").AddOperation(e.NewOperation(1).Add("name", "Boffey")).AddOperation(e.NewOperation(2).Add("name", "BROWN")).AddOperation(e.NewOperation(3).Add("name", "bailey")).AddOperation(e.NewOperation(4).Add("name", "Böhm")).Post()
	t.Sleep(1 * t.Second)
	c.Search("my_index", "user").AddParam("sort", "name.sort").Get()

	// the same field can support multiple sort orders using multi-field for each language and creating the corresponding analyzer for each of these collation.
	c.Mapping("my_index", "_user").AddField("name", e.Dict{"type": "string", "fields": e.Dict{"default": e.Dict{"type": "string", "analyzer": "ducet"}, "french": e.Dict{"type": "string", "analyzer": "french"}, "german": e.Dict{"type": "string", "analyzer": "german_phonebook"}, "swedish": e.Dict{"type": "string", "analyzer": "swedish"}}}).Put()

}
