package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap18 runs example queries from chapter 18 of Elasticsearch the Definitive Guide
// It's an introduction to search in languages using the corresponding analyzers
func chap18() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// specify direcly the language analyzer in field mapping
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("blog", e.NewMapping().AddField("title", e.Dict{e.TYPE: "string", e.ANALYZER: "english"})).Put()
	t.Sleep(1 * t.Second)
	// analyze some text for this field
	c.Analyze("my_index").Field("title").Get("I'm not happy about foxes")

	// double analyze title field: standard analyzer for 'title', english analyzer for subfield 'title.english'
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("blog", e.NewMapping().AddField("title", e.Dict{"type": "string", "fields": e.Dict{"english": e.Dict{"type": "string", "analyzer": "english"}}})).Put()
	// insert some data
	c.Insert("my_index", "blog").Document(1, e.Dict{"title": "I'm happy for this fox"}).Put()
	c.Insert("my_index", "blog").Document(2, e.Dict{"title": "I'm not happy about my fox problem"}).Put()
	t.Sleep(1 * t.Second)
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("type", "most_fields").Add("query", "not happy foxes").AddMultiple("fields", "title", "title.english"))).Pretty().Get()

	// configuring language analyzers,
	// e.g. use english as base analyzer then customize it
	c.Index("my_index").Delete()
	c.Index("my_index").AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_english", e.Dict{"type": "english", e.StemExclusion: []string{"organization", "organizations"}, e.Stopwords: []string{"a", "an", "and", "are", "as", "at", "be", "but", "by", "for", "if", "in", "into", "is", "it", "of", "on", "or", "such", "that", "the", "their", "then", "there", "these", "they", "this", "to", "was", "will", "with"}})).Put()
	// analyze a text with the new analyzer
	t.Sleep(1 * t.Second)
	c.Analyze("my_index").Pretty().Analyzer("my_english").Get("The World Health Organization does not sell organs.")

	// for documents with predominant language use specific index for each one
	// e.g. index for english-specific documents
	c.Index("blogs-en").Mappings("post", e.NewMapping().AddField("title", e.Dict{"type": "string", "fields": e.Dict{"stemmed": e.Dict{"type": "string", "analyzer": "english"}}})).Put()
	// e.g. index for french-specific documents
	c.Index("blogs-fr").Mappings("post", e.NewMapping().AddField("title", e.Dict{"type": "string", "fields": e.Dict{"stemmed": e.Dict{"type": "string", "analyzer": "french"}}})).Put()

	t.Sleep(1 * t.Second)

	// as there is an index for each language, we can specify a preference for particuar languages with 'indices_boost'
	c.Search("blogs-*", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "deja vu").AddMultiple("fields", "title", "title.stemmed").Add("type", "most_fields"))).AddQuery(e.NewQuery(e.IndicesBoost).Add("blogs-en", 3).Add("blogs-fr", 2)).Pretty().Get()
	// it is also possible that a document contains multiple translations for a given field (e.g. title, title_fr, title_es)
	c.Index("movies").Delete()
	c.Index("movies").Mappings("movie", e.NewMapping().AddProperty("title", "type", "string").AddField("title_br", e.Dict{"type": "string", "analyzer": "brazilian"}).AddField("title_cz", e.Dict{"type": "string", "analyzer": "czech"}).AddField("title_en", e.Dict{"type": "string", "analyzer": "english"}).AddField("title_es", e.Dict{"type": "string", "analyzer": "spanish"})).Put()

	// using n-grams
	// e.g. define an custom trigrams analyzer (check "Ngrams for compound words" p. 269)
	c.Index("movies").Delete()
	c.Index("movies").AddAnalyzer(e.NewAnalyzer("filter").Add2("trigrams_filter", e.Dict{"type": "ngram", "min_gram": 3, "max_gram": 3})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("trigrams", e.Dict{"type": "custom", "tokenizer": "standard", "filter": []string{"lowercase", "trigrams_filter"}})).Mappings("movie", e.NewMapping().AddField("title", e.Dict{"type": "string", "fields": e.Dict{"de": e.Dict{"type": "string", "analyzer": "german"}, "en": e.Dict{"type": "string", "analyzer": "english"}, "fr": e.Dict{"type": "string", "analyzer": "french"}, "es": e.Dict{"type": "string", "analyzer": "spanish"}, "general": e.Dict{"type": "string", "analyzer": "trigrams"}}})).Put()
	// e.g. of 'most_fields' search query
	c.Search("movies", "movie").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "club de la lucha").AddMultiple("fields", "title*^1.5", "title.general").Add("type", "most_fields").Add(e.MinimumShouldMatch, "75%"))).Get()
}
