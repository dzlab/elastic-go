package elastic

import ()

// Analyze a structure representing an Elasticsearch query for the Analyze API
type Analyze struct {
	client *Elasticsearch
	parser Parser
	url    string
	params map[string]string
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
		params: make(map[string]string),
	}
}

// Pretty pretiffies the response result
func (analyzer *Analyze) Pretty() *Analyze {
	analyzer.params["pretty"] = ""
	return analyzer
}

// Field adds a field to an anlyze request
func (analyzer *Analyze) Field(field string) *Analyze {
	analyzer.params["field"] = field
	return analyzer
}

// Analyzer adds a named standard Elasticsearch analyzer to the Analyze query
func (analyzer *Analyze) Analyzer(name string) *Analyze {
	analyzer.params["analyzer"] = name
	return analyzer
}

// Get submits an Analyze query to Elasticsearch
// GET /:index/_analyze?field=field_name
func (analyzer *Analyze) Get(body string) {
	// construct the url
	url := urlString(analyzer.url, analyzer.params)

	// construct the body
	analyzer.client.Execute("GET", url, body, analyzer.parser)
}
