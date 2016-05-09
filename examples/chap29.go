package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap29 runs example queries from chapter 29 of Elasticsearch the Definitive Guide.
// It's about dealing with time in aggregations
func chap29() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// use date hisogram to answer how much cars were sold each month?
	// build a bucket each month, format is used to pretty the bucket keys (which will be key_as_string)
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("sales").AddDict(e.DateHistogram, e.Dict{e.Field: "sold", e.Interval: "month", e.Format: "yyyy-MM-dd"})).Get()
	// the previous query will create buckets for interval where 'sold' is greater than zero, and for the rest there will be no values.
	// to get all possible buckets
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("sales").AddDict(e.DateHistogram, e.Dict{e.Field: "sold", e.Interval: "month", e.Format: "yyyy-MM-dd", e.MinDocCount: 0, e.ExtendedBound: e.Dict{e.Min: "2014-01-01", e.Max: "2014-12-31"}})).Get()

	// extended example
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("sales").AddDict(e.DateHistogram, e.Dict{e.Field: "sold", e.Interval: "quarter", e.MinDocCount: 0, e.ExtendedBound: e.Dict{e.Min: "2014-01-01", e.Max: "2014-12-31"}}).AddBucket(e.NewBucket("per_make_sum").AddTerm(e.Field, "make").AddBucket(e.NewBucket("sum_price").AddMetric(e.Sum, e.Field, "price")).AddBucket(e.NewBucket("total_sum").AddMetric(e.Sum, e.Field, "price")))).Get()
}
