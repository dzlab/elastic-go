package elastic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Analyze struct {
	url      string
	field    string
	analyzer string
}

const (
	ANALYZE = "analyze"
)

func (client *Elasticsearch) Analyze(index string) *Analyze {
	var url string = client.request(index, "", -1, ANALYZE)
	return &Analyze{url: url}
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
	url := analyzer.url
	if analyzer.field != "" {
		url = fmt.Sprintf("%s?field=%s", url, analyzer.field)
	} else if analyzer.analyzer != "" {
		url = fmt.Sprintf("%s?analyzer=%s", url, analyzer.analyzer)
	}
	// construct the body
	var data io.Reader = nil
	if body != "" {
		data = bytes.NewReader([]byte(body))
	}
	// submit url
	log.Println("GET", url)
	reader, err := exec("GET", url, data)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
