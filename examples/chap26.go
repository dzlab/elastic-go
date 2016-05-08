package main

import (
	e "github.com/dzlab/elastic-go"
)

// chap26 runs example queries form chapter 26 of Elasticsearch the Definitive Guide.
// It's about learning aggregations by example
func chap26() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// insert some data
	op1 := e.NewOperation(1).Add("price", 10000).Add("color", "red").Add("make", "honda").Add("sold", "2014-10-28")
	op2 := e.NewOperation(2).Add("price", 20000).Add("color", "red").Add("make", "honda").Add("sold", "2014-11-05")
	op3 := e.NewOperation(3).Add("price", 30000).Add("color", "green").Add("make", "ford").Add("sold", "2014-05-18")
	op4 := e.NewOperation(4).Add("price", 15000).Add("color", "blue").Add("make", "toyota").Add("sold", "2014-07-02")
	op5 := e.NewOperation(5).Add("price", 12000).Add("color", "green").Add("make", "toyota").Add("sold", "2014-08-19")
	op6 := e.NewOperation(6).Add("price", 20000).Add("color", "red").Add("make", "honda").Add("sold", "2014-11-05")
	op7 := e.NewOperation(7).Add("price", 80000).Add("color", "red").Add("make", "bmw").Add("sold", "2014-01-01")
	op8 := e.NewOperation(8).Add("price", 25000).Add("color", "blue").Add("make", "ford").Add("sold", "2014-02-12")
	c.Bulk("cars", "transactions").AddOperation(op1).AddOperation(op2).AddOperation(op3).AddOperation(op4).AddOperation(op5).AddOperation(op6).AddOperation(op7).AddOperation(op8).Post()
	// submit an aggregation query
	c.Aggregation("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddTerm("field", "color")).Get()
	// add an addition metric inside bucket: average price
	c.Aggregation("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddTerm("field", "color").AddBucket(e.NewBucket("avg_price").AddMetric(e.Avg, "field", "price"))).Get()
	// bucket inside bucket: another terms bucket that will generated as mush as there is values for field 'make'
	c.Aggregation("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddTerm("field", "color").AddBucket(e.NewBucket("avg_price").AddMetric(e.Avg, "field", "price")).AddBucket(e.NewBucket("make").AddTerm("field", "make"))).Get()
	// add two metics to calculate min/max price for each make
	c.Aggregation("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("colors").AddTerm("field", "color").AddBucket(e.NewBucket("avg_price").AddMetric(e.Avg, "field", "price")).AddBucket(e.NewBucket("make").AddTerm("field", "make").AddBucket(e.NewBucket("min_price").AddMetric(e.Min, "price")).AddBucket(e.NewBucket("max_price").AddMetric(e.Max, "price")))).Get()
}
