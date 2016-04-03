package elastic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

/*
 * Elasticsearch client
 */
type Elasticsearch struct {
	Addr string
}

/*
 * Build the url of an API request call
 */
func (this *Elasticsearch) request(index, class string, id int64, request string) string {
	var url string
	if index == "" {
		url = fmt.Sprintf("http://%s/_%s", this.Addr, request)
	} else if class == "" {
		url = fmt.Sprintf("http://%s/%s/_%s", this.Addr, index, request)
	} else if id < 0 {
		url = fmt.Sprintf("http://%s/%s/%s/_%s", this.Addr, index, class, request)
	} else {
		url = fmt.Sprintf("http://%s/%s/%s/%d/_%s", this.Addr, index, class, id, request)
	}
	return url
}

/*
 * Elasticsearch failure representation
 * e.g.:{"error":{"root_cause":[{"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"}],"type":"no_shard_available_action_exception","reason":"No shard available for [org.elasticsearch.action.admin.indices.analyze.AnalyzeRequest@74508901]"},"status":503}
 * e.g.:{"error":{"root_cause":[{"type":"index_already_exists_exception","reason":"already exists","index":"my_index"}],"type":"index_already_exists_exception","reason":"already exists","index":"my_index"},"status":400}
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
func String(obj interface{}) string {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
	}
	return string(marshaled)
}

/*
 * Construct a url
 */
func urlString(prefix string, params map[string]string) string {
	url := prefix
	if len(params) > 0 {
		if strings.Contains(url, "?") {
			if len(params) > 0 {
				url += "&"
			}
		} else {
			url += "?"
		}
		for name, value := range params {
			url += name
			if value != "" {
				url += "=" + value
			}
			url += "&"
		}
		url = url[:len(url)-1]
	}
	return url
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
