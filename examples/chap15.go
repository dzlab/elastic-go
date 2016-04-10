package main

import (
	e "github.com/dzlab/elastic-go"
	"time"
)

/*
 * Examples of queries based on Elasticsearch Definitive Guide, chapter 15
 * Proximity matching search examples
 */
func chap15() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// phrase matching with 'match_phrase' search query
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().Add("title", "quick brown fox"))).Get()
	// a rewrite of the previsous query using 'match' query
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("title").Add("query", "quick brown fox").Add("type", "phrase")))).Get()

	// 'match_phrase' query can use terms position to search for document, this position can be seen with 'Analyze' query
	c.Analyze("").Analyzer("standard").Get("quick brown fox")

	// we can introduce flexibity in the search query by using 'slop' parameters
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match_phrase").AddQuery(e.NewQuery("title").Add("query", "quick fox").Add("slop", 1)))).Get()

	// multi-value fields react surprisingly to 'match_phrase' queries
	c.Insert("my_index", "groups").Document(1, e.Dict{"names": []string{"John Abraham", "Lincoln Smith"}}).Put()
	time.Sleep(1 * time.Second)
	c.Search("my_index", "groups").AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().Add("names", "Abraham Lincoln"))).Get()
	// to avoid successive documents to appear in search result, use 'position_offset_gap' when creating the index in order to increase offset between these documents
	c.Index("my_type/groups").Delete()
	c.Mapping("my_type", "groups").AddProperty("names", "type", "string").AddProperty("names", e.POSITION_OFFSET_GAP, 100).Put()

	// proximity query (phrase query with 'slop' higher than 0) includes proximity of query terms in the result '_score' field
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("title").Add("query", "quick dog").Add("slop", 50)))).Get()

	// proximity queries can be compbined with 'match' query to filter irrelevant documnets
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddMust(e.NewMatch().AddQuery(e.NewQuery("title").Add("query", "quick brown fox").Add(e.MINIMUM_SHOULD_MATCH, "30%"))).AddShould(e.NewMatchPhrase().AddQuery(e.NewQuery("title").Add("query", "quick brwon fox").Add(e.SLOP, 50))))).Get()

	// beware of performance overhead, as a simple 'term' query is 10 times as fast as a 'phrase' query, and 20 times as fast as a proximity query (phrase query with 'slop')
	// to increase performance, one option will be to reduce number of documents
	// we case use 'match' query, to catch relevant documents than rescore using some scoring algorithm
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().AddQuery(e.NewQuery("title").Add("query", "quick brown fox").Add(e.MINIMUM_SHOULD_MATCH, "30%")))).AddQuery(e.NewRescore().Add(e.WINDOW_SIZE, 50).AddQuery(e.NewQuery("query").AddQuery(e.NewRescoreQuery().AddQuery(e.NewMatchPhrase().AddQuery(e.NewQuery("title").Add("query", "quick brown fox").Add("slop", 50)))))).Get()

	// instead of indexing words separately, we can index bigrams (or shingles) to retain more of the context in which words occured
	// producing shingles
	c.Index("my_index").Delete()
	c.Index("my_index").SetShardsNb(1).AddAnalyzer(e.NewAnalyzer("filter").Add2("my_shingle_filter", e.Dict{e.TYPE: "string", e.MIN_SHINGLE_SIZE: 2, e.MAX_SHINGLE_SIZE: 2, e.OUTPUT_UNIGRAMS: false})).AddAnalyzer(e.NewAnalyzer("analyzer").Add2("my_shingle_analyzer", e.Dict{e.TYPE: "custom", e.TOKENIZER: "standard", e.FILTER: []string{"lowercase", "my_shingle_filter"}})).Put()

	// index unigrams and bigrams separately
	// create title field as multi-field: unigrams(title), (title.shingles)
	c.Mapping("my_index", "my_type").AddDocumentType(e.NewDocType("my_type").AddTemplate(e.NewTemplate("properties").AddProperty("title", e.Dict{e.TYPE: "string", "fields": e.Dict{"shingles": e.Dict{e.TYPE: "string", "analyzer": "my_shingle_analyzer"}}}))).Put()
	// insert some documents
	c.Bulk("my_index", "my_type").AddOperation(e.NewOperation(1).Add("title", "Sue ate the alligator")).AddOperation(e.NewOperation(2).Add("title", "The aligator ate Sue")).AddOperation(e.NewOperation(3).Add("title", "Sue never goes anywhere without her alligator skin purse")).Post()

	// searching for shingles
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("title", "the hungry alligator ate sue"))).Get()
	// let's add 'shingles' to act as signal and increase relevance score
	c.Search("my_index", "my_type").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddMust(e.NewMatch().Add("title", "the hungry alligator ate sue")).AddShould(e.NewMatch().Add("title.shingles", "the hungry alligator ate sue")))).Get()
}
