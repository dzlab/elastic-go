package elastic

import (
	"fmt"
)

type Analyze struct {
	client   *Elasticsearch
	parser   Parser
	url      string
	field    string
	analyzer string
}

const (
	ANALYZE = "analyze"
)

func (client *Elasticsearch) Analyze(index string) *Analyze {
	var url string = client.request(index, "", -1, ANALYZE)
	return &Analyze{
		client: client,
		parser: &AnalyzeResultParser{},
		url:    url,
	}
}

/*
 * Analyze a field
 */
func (analyzer *Analyze) Field(field string) *Analyze {
	analyzer.field = field
	return analyzer
}

/*
 * Analyze an analyzer given by name
 */
func (analyzer *Analyze) Analyzer(name string) *Analyze {
	analyzer.analyzer = name
	return analyzer
}

/*
 * GET /:index/_analyze?field=field_name
 */
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
