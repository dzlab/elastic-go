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

}
