package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap31 runs example queries from chapter 31 of Elasticsearch the Definitive Guide.
// It's about using filtering queries along with aggregations
func chap31() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// Find all cars over $10000 and calculate the average price of these cars
	c.Aggs("cars", "transactions").SetMetric(e.Count).AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("filter").AddQuery(e.NewQuery("rage").AddQuery(e.NewQuery("price").Add("gte", "10000"))))).AddBucket(e.NewBucket("single_avg_price").AddMetric(e.Avg, e.Field, "price")).Get()

	// filtering aggregation results
	c.Aggs("cars", "transactions").SetMetric(e.Count).AddQuery(e.NewMatch().Add("make", "ford")).AddBucket(e.NewBucket("recent_sales").AddDict(e.FilterBucket, e.Dict{"range": e.Dict{"sold": e.Dict{"from": "now-1M"}}}).AddBucket(e.NewBucket("average_price").AddMetric(e.Avg, e.Field, "price"))).Get()

	// filtering results without afecting the query scope.
	c.Aggs("cars", "transactions").SetMetric(e.Count).AddQuery(e.NewMatch().Add("make", "ford")).AddPostFilter(e.NewTerm().Add("color", "green")).AddBucket(e.NewBucket("all_colors").AddTerm(e.Field, "color")).Get()
}
