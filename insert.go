package elastic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

/*
 * a request representing an insert
 */
type Insert struct {
	url string
	id  int64
	doc interface{}
}

/*
 * Create an Insert request, that will submit a new document to elastic search
 */
func (this *Elasticsearch) Insert(index, doctype string) *Insert {
	var url string = fmt.Sprintf("http://%s/%s/%s", this.Addr, index, doctype)
	return &Insert{url: url}
}

/*
 * Create a new Insert query (for test purpose)
 */
func newInsert() *Insert {
	return &Insert{}
}

/*
 * Set the document to insert
 */
func (this *Insert) Document(id int64, doc interface{}) *Insert {
	this.id = id
	this.doc = doc
	return this
}

/*
 * Get a string representation of the document
 */
func (this *Insert) String() string {
	return String(this.doc)
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * PUT /:index/:type/:id
 */
func (this *Insert) Put() {
	// construct the url
	url := fmt.Sprintf("%s/%d", this.url, this.id)
	// construct the body
	query := this.String()
	var body io.Reader
	if query != "" {
		body = bytes.NewReader([]byte(query))
	}
	// submit the request
	log.Println("PUT", url, query)
	reader, err := exec("PUT", url, body)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
