package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap34 examples from chapter 34 of Elasticsearch the Definitive Guide. It's about Significant terms.
func chap34() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// execute a search without query to see random sampling
	c.Search("mlmovies", "").Get()
	// look at mlratings
	c.Search("mlratings", "").Get()

	// recommendation based on popularity
	c.Search("mlmovies", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("title", "Talladega Nights"))).Get()
	// with the ID, filter ratings and apply a terms aggregation to find most popular .. from people that also rated ..
	c.Aggs("mlratings", "").SetMetric(e.Count).AddQuery(e.NewQuery("filtered").AddQuery(e.NewFilter().AddQuery(e.NewTerm().Add("movie", 46970)))).Add(e.NewBucket("most_popular").AddDict(e.Terms, e.Dict{e.Field: "movie", e.Size: 6})).Get()
	// we need to colorate this with their original titles
	c.Search("mlmovies", "").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewFilter().AddQuery(e.NewQuery("ids").AddMultiple("values", 2571, 318, 296, 2959, 260))))).Get()
	// this a recommendation of most popular .. it is not a recommendation based on .., we can verify this by removing the filter part and comparing the results.
	c.Aggs("mlratings", "").SetMetric(e.Count).Add(e.NewBucket("most_popular").AddDict(e.Terms, e.Dict{e.Field: "movie", e.Size: 5})).Get()
	// just checking the most popular .. is not sufficient to build good discriminating recommender
	c.Aggs("mlratings", "").SetMetric(e.Count).AddQuery(e.NewQuery("filtered").AddQuery(e.NewFilter().AddQuery(e.NewTerm().Add("movie", 46970)))).Add(e.NewBucket("most_sig").AddDict(e.SignificantTerms, e.Dict{e.Field: "movie", e.Size: 6})).Get()
}
