package elastic

import (
	"bytes"
	"encoding/json"
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
 * Elasticsearch failure representation
 * e.g.:{"error":{"root_cause":[{"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"}],"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"},"status":503}
 */
type Failure struct {
	kind   string `type`
	reason string `reason`
	status int    `json`
}

/*
 * Elasticsearch success representation
 * e.g.: {"acknowledged":true}
 */
type Success struct {
	acknowledged bool
}

/*
 * Elasticsearch unvalid representation
 * e.g.: {"valid":false,"_shards":{"total":1,"successful":1,"failed":0},"explanations":[{"index":"gb","valid":false,"error":"org.elasticsearch.index.query.QueryParsingException: No query registered for [tweet]"}]}
 */
type Unvalid struct {
	valid       bool
	shards      Dict
	explanation []Dict
}

/*
 * Return a string representation of the dictionary
 */
func String(dict Dict) string {
	marshaled, err := json.Marshal(dict)
	if err != nil {
		log.Println(err)
	}
	return string(marshaled)
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
