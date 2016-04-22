package elastic

import (
	"fmt"
)

// Analyze a structure representing an Elasticsearch query for the Analyze API
type Analyze struct {
	client   *Elasticsearch
	parser   Parser
	url      string
	field    string
	analyzer string
}

const (
	// ANALYZE a constant for Analyze query name
	ANALYZE = "analyze"
)

// Analyze returns an new Analyze request on the given index
func (client *Elasticsearch) Analyze(index string) *Analyze {
	url := client.request(index, "", -1, ANALYZE)
	return &Analyze{
		client: client,
		parser: &AnalyzeResultParser{},
		url:    url,
	}
}

// Field adds a field to an anlyze request
func (analyzer *Analyze) Field(field string) *Analyze {
	analyzer.field = field
	return analyzer
}

// Analyzer adds a named standard Elasticsearch analyzer to the Analyze query
func (analyzer *Analyze) Analyzer(name string) *Analyze {
	analyzer.analyzer = name
	return analyzer
}

// Get submits an Analyze query to Elasticsearch
// GET /:index/_analyze?field=field_name
func (analyzer *Analyze) Get(body string) {
	// construct the url
	url := analyzer.url
	if analyzer.field != "" {
		url = fmt.Sprintf("%s?field=%s", url, analyzer.field)
	} else if analyzer.analyzer != "" {
		url = fmt.Sprintf("%s?analyzer=%s", url, analyzer.analyzer)
	}
	// construct the body
	analyzer.client.Execute("GET", url, body, analyzer.parser)
}
