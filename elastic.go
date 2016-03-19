package elastic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/*
 * Elasticsearch client
 */
type Elasticsearch struct {
	Addr string
}

/*
 * mappings between the json fields and how Elasticsearch store them
 */
type Mapping struct {
	url string
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/_mapping/:type
 */
func (this *Elasticsearch) Mapping(index, class string) *Mapping {
	url := fmt.Sprintf("http://%s/%s/_mapping/%s", this.Addr, index, class)
	return &Mapping{url: url}
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/_mapping/:type
 */
func (this *Mapping) Get() {
	reader, err := exec("GET", this.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

/*
 * Update a mappings between the json fields and how Elasticsearch store them
 * PUT /:index/_mapping/:type
 */
func (this *Mapping) Put(body string) {
	reader, err := exec("PUT", this.url, bytes.NewReader([]byte(body)))
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

type Analyze struct {
	url string
}

func (this *Elasticsearch) Analyze(index string) *Analyze {
	url := fmt.Sprintf("http://%s/%s/_analyze", this.Addr, index)
	return &Analyze{url: url}
}

/*
 * GET /:index/_analyze?field=field_name
 */
func (this *Analyze) Get(field string) {
	url := fmt.Sprintf("%s?field=%s", this.url, field)
	log.Println(url)
	reader, err := exec("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

/*
 * Execute a REST request
 */
func exec(method, url string, body io.Reader) (io.Reader, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
