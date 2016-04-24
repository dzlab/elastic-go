package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

// chap17 runs example queries from chapter 17 of Elasticsearch the Definitive Guide
// It's about controlling relevance score of a given document
// relevance score is based on Term Frequency/Inverse Document Frequency and Vector Space Model, in addition to a coordination factor, field length normalization and term/query clause boosting
func chap17() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// disable term frequency in relevance calculation when we don't care about term weight and/or position in matching documents
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("doc", e.NewMapping().AddField("text", e.Dict{e.TYPE: "string", e.IndexOptions: "docs"})).Put()

	// disable field length normalisation in calculating relevance score
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("doc", e.NewMapping().AddField("text", e.Dict{e.TYPE: "string", e.Norms: e.Dict{"enabled": false}})).Put()

	// check the explanation of a search query to see the scoring factors in action
	c.Index("my_index").Delete()
	c.Insert("my_index", "doc").Document(1, e.Dict{"text": "quick brown fox"}).Put()
	t.Sleep(1 * t.Second)
	c.Search("my_index", "doc").Pretty().AddParam("explain", "").AddQuery(e.NewQuery("query").AddQuery(e.NewTerm().Add("text", "fox"))).Get()

	// disable query coordination function in a query with synonyms
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().Add(e.DisableCoord, true).AddShould(e.NewTerm().Add("text", "jump")).AddShould(e.NewTerm().Add("text", "hop")).AddShould(e.NewTerm().Add("text", "leap")))).Get()

	// boosting an index
	c.Search("docs_2014_*", "").AddQuery(e.NewQuery(e.IndicesBoost).Add("docs_2014_10", 3).Add("docs_2014_09", 2)).AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("text", "quick brown fox"))).Get()

	// boosting query that downgrade document about apple the fruit
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBoosting().SetNegativeBoost(0.5).AddPositive("match", e.Dict{"text": "apple"}).AddNegative("match", e.Dict{"text": "pie tart fruit crumble tree"}))).Get()

	// constant score query that assigns a score of 1 to any document
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewConstantScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "wifi")))).AddShould(e.NewConstantScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "garden")))).AddShould(e.NewConstantScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "pool")))))).Get()
	// a specific boost value to a clause
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewConstantScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "wifi")))).AddShould(e.NewConstantScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "garden")))).AddShould(e.NewConstantScore().Add(e.Boost, 2).AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("description", "pool")))))).Get()

	// full text search with boosting (more relevance) based on popularity
	c.Insert("blogposts", "post").Document(1, e.Dict{"title": "About popularity", "content": "In this post we will talk about...", "votes": 6}).Put()
	c.Search("blogposts", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "popularity").Add("fields", []string{"title", "content"}))).AddQuery(e.NewQuery(e.FieldValueFactor).Add("field", "votes")))).Get()
	// a better way to incorporate popularity is by using a modifier (e.g. log1p) so that first few votes count a lot, but subsequent votes less
	c.Search("blogposts", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "popularity").Add("fields", []string{"title", "content"}))).AddQuery(e.NewQuery(e.FieldValueFactor).Add("field", "votes").Add("modifier", "log1p")))).Get()
	// use in addition a factor
	c.Search("blogposts", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "popularity").Add("fields", []string{"title", "content"}))).AddQuery(e.NewQuery(e.FieldValueFactor).Add("field", "votes").Add(e.Modifer, "log1p").Add(e.Factor, 2)))).Get()
	// use boost_mode to modifiy how calculated score is combined with _score
	c.Search("blogposts", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "popularity").Add("fields", []string{"title", "content"}))).AddQuery(e.NewQuery(e.FieldValueFactor).Add("field", "votes").Add(e.Modifer, "log1p").Add("factor", 0.1)).Add(e.BoostMode, "sum"))).Get()
	// cap the maximum of the scoring function
	c.Search("blogposts", "post").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "popularity").Add("fields", []string{"title", "content"}))).AddQuery(e.NewQuery(e.FieldValueFactor).Add("field", "votes").Add(e.Modifer, "log1p").Add("factor", 0.1)).Add(e.BoostMode, "sum").Add(e.MaxBoost, 1.5))).Get()

	// boosting filtered subsets
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewFilter().Add("term", e.Dict{"city": "Barcelona"})).AddMultiple("functions", e.Dict{e.Filter: e.NewQuery("term").Add("features", "wifi").Dict(), e.Weight: 1}, e.Dict{e.Filter: e.NewQuery("term").Add("features", "garden").Dict(), e.Weight: 1}, e.Dict{e.Filter: e.NewQuery("term").Add("features", "pool").Dict(), e.Weight: 2}).Add(e.ScoreMode, "sum"))).Get()
	// introduce some randomness so that documents with similar score get same exposuer with same order for each user (i.e. consistently random) in the seed parameter
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddQuery(e.NewFilter().AddQuery(e.NewTerm().Add("city", "Barcelona"))).AddMultiple("functions", e.Dict{e.Filter: e.NewTerm().Add("features", "wifi").Dict(), e.Weight: 1}, e.Dict{e.Filter: e.NewTerm().Add("features", "garden").Dict(), e.Weight: 1}, e.Dict{e.Filter: e.NewTerm().Add("features", "pool").Dict(), e.Weight: 2}, e.NewQuery(e.RandomScore).Add("seed", "the users session id").Dict()).Add(e.ScoreMode, "sum"))).Get()

	// decay function: the closer the better
	// e.g. find a place to rent near center of london and not exceeding 100£ the night
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddMultiple("functions", e.NewQuery("gauss").AddQuery(e.NewQuery("location").Add("origin", e.Dict{"lat": 51.5, "lon": 0.12}).Add("offset", "2km").Add("scale", "3km")).Dict(), e.NewQuery("gauss").AddQuery(e.NewQuery("price").Add("origin", "50").Add("offset", 50).Add("scale", "20")).Add(e.Weight, 2).Dict()))).Pretty().Get()
	// use a custom Groovy script to score documents
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewFunctionScore().AddMultiple("functions", e.NewQuery("gauss").AddQuery(e.NewQuery("location").Add("origin", e.Dict{"lat": 51.5, "lon": 0.12}).Add("offset", "2km").Add("scale", "3km")).Dict(), e.NewQuery("gauss").AddQuery(e.NewQuery("price").Add("origin", "50").Add("offset", 50).Add("scale", "20")).Add(e.Weight, 2).Dict(), e.NewQuery(e.ScriptScore).AddQuery(e.NewQuery("params").Add("threshold", 80).Add("discount", 0.1).Add("target", 10)).Add("script", "price = doc['price'].value; margin=doc['margin'].value;if(price<threshold){return price * margin/target}; return price * (1-discount)*margin/target").Dict()))).Pretty().Get()

	// changing similarities
	c.Index("my_index").Mappings("doc", e.NewMapping().AddField("title", e.Dict{e.TYPE: "string", e.Similarity: "BM25"}).AddField("body", e.Dict{e.TYPE: "string", e.Similarity: "default"})).Put()
	// configuring BM25, e.g. disable field length normalization
	c.Index("my_index").Settings(e.NewQuery("similarity").AddQuery(e.NewQuery("my_bm25").Add("type", "BM25").Add("b", 0)).Dict()).Mappings("doc", e.NewMapping().AddField("title", e.Dict{e.TYPE: "string", e.Similarity: "my_bm25"}).AddField("body", e.Dict{e.TYPE: "string", e.Similarity: "BM25"})).Put()
}
