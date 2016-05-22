package main

import (
	e "github.com/dzlab/elastic-go"
	t "time"
)

func chap33() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// use cardinality to approximate count of uniquer values
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("distinct_colors").AddMetric(e.Cardiality, e.Field, "color")).Get()
	// a more interesting question: how many colors were sold each month
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("months").AddDict(e.DateHistogram, e.Dict{e.Field: "sold", e.Interval: "month"}).AddBucket(e.NewBucket("distinct_colors").AddMetric(e.Cardiality, e.Field, "color"))).Get()

	// configure the cardinaility HuperLogLog (HLL) algorithm
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("distinct_colors").AddDict(e.Cardiality, e.Dict{e.Field: "color", e.PrecisionThreshold: 100})).Get()

	// calculte hashes for HLL at index time for speedup
	c.Index("cars").Delete()
	c.Index("cars").Mappings("transactions", e.NewMapping().AddField("color", e.Dict{e.Type: "string", "fields": e.Dict{"hash": e.Dict{e.Type: "murmur3"}}})).Put()
	c.Bulk("cars", "transactions").AddOperation(e.NewOperation(1).Add("price", 10000).Add("color", "red").Add("make", "honda").Add("sold", "2014-10-28")).AddOperation(e.NewOperation(2).Add("price", 20000).Add("color", "red").Add("make", "honda").Add("sold", "2014-11-05")).AddOperation(e.NewOperation(3).Add("price", 30000).Add("color", "green").Add("make", "ford").Add("sold", "2014-05-18")).AddOperation(e.NewOperation(4).Add("price", 15000).Add("color", "blue").Add("make", "toyota").Add("sold", "2014-07-02")).AddOperation(e.NewOperation(5).Add("price", 12000).Add("color", "green").Add("make", "toyota").Add("sold", "2014-08-19")).AddOperation(e.NewOperation(6).Add("price", 20000).Add("color", "red").Add("make", "honda").Add("sold", "2014-11-05")).AddOperation(e.NewOperation(7).Add("price", 80000).Add("color", "red").Add("make", "bmw").Add("sold", "2014-01-01")).AddOperation(e.NewOperation(8).Add("price", 25000).Add("color", "blue").Add("make", "ford").Add("sold", "2014-02-12")).Post()
	// cardinality aggregtion on the hashed field
	c.Aggs("cars", "transactions").SetMetric(e.Count).Add(e.NewBucket("distinct_colors").AddMetric(e.Cardiality, e.Field, "color.hash")).Get()

	// percentile metric
	c.Index("website").Delete()
	c.Bulk("website", "logs").AddOperation(e.NewOperation(1).Add("latency", 100).Add("zone", "US").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(2).Add("latency", 80).Add("zone", "US").Add("timestamp", "2014-10-29")).AddOperation(e.NewOperation(3).Add("latency", 99).Add("zone", "US").Add("timestamp", "2014-10-29")).AddOperation(e.NewOperation(4).Add("latency", 102).Add("zone", "US").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(5).Add("latency", 75).Add("zone", "US").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(6).Add("latency", 82).Add("zone", "US").Add("timestamp", "2014-10-29")).AddOperation(e.NewOperation(7).Add("latency", 100).Add("zone", "EU").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(8).Add("latency", 280).Add("zone", "EU").Add("timestamp", "2014-10-29")).AddOperation(e.NewOperation(9).Add("latency", 155).Add("zone", "EU").Add("timestamp", "2014-10-29")).AddOperation(e.NewOperation(10).Add("latency", 623).Add("zone", "EU").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(11).Add("latency", 380).Add("zone", "EU").Add("timestamp", "2014-10-28")).AddOperation(e.NewOperation(12).Add("latency", 319).Add("zone", "EU").Add("timestamp", "2014-10-29")).Post()

	t.Sleep(1 * t.Second)

	// run percentile to get an array of predefined percentiles
	c.Aggs("website", "logs").SetMetric(e.Count).Add(e.NewBucket("load_times").AddMetric(e.Percentiles, e.Field, "latency")).Add(e.NewBucket("avg_load_time").AddMetric(e.Avg, e.Field, "latency")).Get()
	// let see if this is corelated with the geographic locations
	c.Aggs("website", "logs").SetMetric(e.Count).Add(e.NewBucket("zones").AddTerm(e.Field, "zone").AddBucket(e.NewBucket("load_times").AddDict(e.Percentiles, e.Dict{e.Field: "latency", e.Percents: []float32{50, 95.0, 99.0}})).AddBucket(e.NewBucket("load_avg").AddMetric(e.Avg, e.Field, "latency"))).Get()

	// use percentile_metric to find the percentile corresponding to 200 ms
	c.Aggs("website", "logs").SetMetric(e.Count).Add(e.NewBucket("zones").AddTerm(e.Field, "zone").AddBucket(e.NewBucket("load_times").AddDict(e.PercentileRanks, e.Dict{e.Field: "latency", e.Values: []int{210, 800}}))).Get()
}
