package elastic

import ()

const (
	// Aggs abreviateed constant name for the Aggregation query.
	Aggs = "aggs"
	// Aggregations constant name for the Aggregation query.
	Aggregations = "aggregations"
	// Terms constant name of terms Bucket
	Terms = "terms"
)

// Constant name of Elasticsearch metrics
const (
	// Count constant name of 'count' metric.
	Count = "count"
	// Avg constant name of 'avg' metric.
	Avg = "avg"
	// Min constant name of 'min' metric.
	Min = "min"
	// Max constant name of 'max' metric.
	Max = "max"
)

type Aggregation struct {
	client *Elasticsearch
	parser *Parser
	url    string
	params map[string]string
	query  Dict
}

func (client *Elasticsearch) Aggs(index, doc string) *Aggregation {
	url := client.request(index, doc, -1, SEARCH)
	return &Aggregation{
		client: client,
		url:    url,
		params: make(map[string]string),
		query:  make(Dict),
	}
}

// urlString constructs the url of this Search API call
func (agg *Aggregation) urlString() string {
	return urlString(agg.url, agg.params)
}

// String returns a string representation of this Search API call
func (agg *Aggregation) String() string {
	body := ""
	if len(agg.query) > 0 {
		dict := make(Dict)
		dict[Aggs] = agg.query
		body = String(dict)
	}
	return body
}

// Get submits request mappings between the json fields and how Elasticsearch store them
// GET /:index/:type/_search
func (agg *Aggregation) Get() {
	// construct the url
	url := agg.urlString()
	// construct the body
	query := agg.String()

	agg.client.Execute("GET", url, query, agg.parser)
}

func (agg *Aggregation) SetMetric(name string) *Aggregation {
	agg.params[SearchType] = name
	return agg
}

type Bucket struct {
	name  string
	query Dict
}

func NewBucket(name) *Bucket {
	return &Bucket{
		name:  name,
		query: make(Dict),
	}
}

func (bucket *Bucket) AddTerm(name string, value interface{}) *Bucket {
	bucket.AddMetric(Terms, name, value)
	return bucket
}

func (bucket *Bucket) AddMetric(metric, name, string, value interface{}) *Bucket {
	if bucket.query[metric] == nil {
		bucket.query[metric] = make(Dict)
	}
	bucket.query[metric].(Dict)[name] = value
	return bucket
}

func (bucket *Bucket) AddBucket(b *Bucket) *Bucket {
	if bucket.query[Aggs] == nil {
		bucket.query[Aggs] = make(Dict)
	}
	bucket.query[Aggs].(Dict)[b.name] = b.query
	return bucket
}

func (agg *Aggregation) Add(bucket Bucket) *Aggregation {
	agg.query[bucket.name] = bucket.query
	return agg
}
