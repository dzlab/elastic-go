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
