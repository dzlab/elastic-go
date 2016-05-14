package main

import (
	e "github.com/dzlab/elastic-go"
)

func chap32() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}
	// intrinsic sort on bucket's document count
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddTerm(e.Field, "color").SetOrder(e.Terms, "_count", "asc")).Get()

	// sort by a single value metric
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddBucket(e.NewBucket("avg_price").AddMetric(e.Avg, e.Field, "price")).AddTerm(e.Field, "color").SetOrder(e.Terms, "avg_price", "asc")).Get()

	// sort by a multiple values metric
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddBucket(e.NewBucket("stats").AddMetric(e.ExtendedStats, e.Field, "price")).AddTerm(e.Field, "color").SetOrder(e.Terms, "stats.variance", "asc")).Get()

	// sory by metric of a deeply nested bucket
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddDict(e.Histogram, e.Dict{e.Field: "price", e.Interval: 20000}).SetOrder(e.Histogram, "red_green_cars>stats.variance", "asc").AddBucket(e.NewBucket("red_green_cars").AddDict(e.FilterBucket, e.Dict{e.Terms: e.Dict{"color": []string{"red", "green"}}}).AddBucket(e.NewBucket("stats").AddMetric(e.ExtendedStats, e.Field, "price")))).Get()
}
