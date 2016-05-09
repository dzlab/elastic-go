package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap28 runs example queries from chapter 28 of Elasticsearch the Definitive Guide.
// It's about building bar charts
func chap28() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// create a histogram bucket of 20000 for the car price, and an nested bucket of size one that gives the bucket with highest count (i.e. top number of cars sold of a maker)
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("price").AddMetric(e.Histogram, e.Field, "price").AddMetric(e.Histogram, e.Interval, 20000).AddBucket(e.NewBucket("make").AddTerm(e.Field, "make").AddTerm(e.Size, 1))).Get()

	// chart bar of popular makes, their average price and standard error.
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("makes").AddTerm(e.Field, "make").AddTerm(e.Size, 10).AddBucket(e.NewBucket("stats").AddMetric(e.ExtendedStats, e.Field, "price"))).Get()
}
