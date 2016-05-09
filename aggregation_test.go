package elastic

import (
	"testing"
)

// newAggs create a new Aggregation for testing
func newAggs() *Aggregation {
	return &Aggregation{
		client: nil,
		url:    "/",
		params: make(map[string]string),
		query:  make(Dict),
	}
}

// test for Aggregations
func TestAggregations(t *testing.T) {
	actual := []string{
		newAggs().Add(NewBucket("colors").AddTerm("field", "color")).String(),
		newAggs().Add(NewBucket("colors").AddTerm("field", "color").AddBucket(NewBucket("avg_price").AddMetric(Avg, "field", "price")).AddBucket(NewBucket("make").AddTerm("field", "make").AddBucket(NewBucket("min_price").AddMetric(Min, Field, "price")).AddBucket(NewBucket("max_price").AddMetric(Max, Field, "price")))).String(),
		newAggs().Add(NewBucket("makes").AddTerm(Field, "make").AddTerm(Size, 10).AddBucket(NewBucket("stats").AddMetric(ExtendedStats, Field, "price"))).String(),
	}
	expected := []string{
		`{"aggs":{"colors":{"terms":{"field":"color"}}}}`,
		`{"aggs":{"colors":{"aggs":{"avg_price":{"avg":{"field":"price"}},"make":{"aggs":{"max_price":{"max":{"field":"price"}},"min_price":{"min":{"field":"price"}}},"terms":{"field":"make"}}},"terms":{"field":"color"}}}}`,
		`{"aggs":{"makes":{"aggs":{"stats":{"extended_stats":{"field":"price"}}},"terms":{"field":"make","size":10}}}}`,
	}
	equals(t, actual, expected)
}

// tests for query scope functionality
func TestQueryScope(t *testing.T) {
	actual := []string{
		newAggs().AddQuery(NewMatch().Add("make", "ford")).Add(NewBucket("colors").AddTerm(Field, "color")).String(),
		newAggs().AddQuery(NewQuery("filtered").AddQuery(NewQuery("filter").AddQuery(NewQuery("rage").AddQuery(NewQuery("price").Add("gte", 10000))))).Add(NewBucket("single_avg_price").AddMetric(Avg, Field, "price")).String(),
	}
	expected := []string{
		`{"aggs":{"colors":{"terms":{"field":"color"}}},"query":{"match":{"make":"ford"}}}`,
		`{"aggs":{"single_avg_price":{"avg":{"field":"price"}}},"query":{"filtered":{"filter":{"rage":{"price":{"gte":10000}}}}}}`,
	}
	equals(t, actual, expected)
}
