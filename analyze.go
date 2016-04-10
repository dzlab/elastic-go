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

func (this *Elasticsearch) Analyze(index string) *Analyze {
	var url string = this.request(index, "", -1, ANALYZE)
	return &Analyze{url: url}
}

/*
 * Analyze a field
 */
func (this *Analyze) Field(field string) *Analyze {
	this.field = field
	return this
}

/*
 * Analyze an analyzer
 */
func (this *Analyze) Analyzer(analyzer string) *Analyze {
	this.analyzer = analyzer
	return this
}

/*
 * GET /:index/_analyze?field=field_name
 */
func (this *Analyze) Get(body string) {
	url := this.url
	if this.field != "" {
		url = fmt.Sprintf("%s?field=%s", url, this.field)
	} else if this.analyzer != "" {
		url = fmt.Sprintf("%s?analyzer=%s", url, this.analyzer)
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
