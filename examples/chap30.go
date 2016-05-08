package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap30 runs example queries from chapter 30 of Elasticsearch the Definitive Guide.
// It's about scoping aggragations, i.e. defining a query to run aggregations on matched documents.
func chap30() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// how many colors are Ford cars are available in?
	c.Aggs("cars", "transactions").AddQuery(e.NewMatch().Add("make", "ford")).AddBucket(e.NewBucket("colors").AddTerm(e.Field, "color")).Get()

	// global bucket to by pass aggregation scope.
	// compare ford cars sales with all cars sales
	c.Aggs("cars", "transactions").SetMetric(e.Count).AddQuery(e.NewMatch().Add("make", "ford")).AddBucket(e.NewBucket("single_avg_price").AddMetric(e.Avg, e.Field, "price")).AddBucket(e.NewBucket("all").AddDict(e.Global, e.Dict{}).AddBucket(e.NewBucket("avg_price").AddMetric(e.Avg, e.Field, "price"))).Get()
}
