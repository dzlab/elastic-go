package elastic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	BULK = "bulk"
)

type Bulk struct {
	url string
	ops []Dict
}

/*
 * A bulk operation
 */
type Operation struct {
	_id int
	doc Dict
}

/*
 * Create a new Bulk of operations
 */
func newBulk() *Bulk {
	return &Bulk{url: "", ops: []Dict{}}
}

/*
 * Create a new operation with the given id
 */
func NewOperation(id int) *Operation {
	return &Operation{_id: id, doc: make(Dict)}
}

/*
 * Add a field to this document
 */
func (this *Operation) Add(name string, value interface{}) *Operation {
	this.doc[name] = value
	return this
}

/*
 * Create a new Bulk operations
 */
func (this *Elasticsearch) Bulk(index, docType string) *Bulk {
	url := this.request(index, docType, -1, BULK)
	return &Bulk{url: url, ops: []Dict{}}
}

/*
 * Add an operation to this bulk
 */
func (this *Bulk) AddOperation(op *Operation) *Bulk {
	indexOp := Dict{"index": Dict{"_id": op._id}}
	this.ops = append(this.ops, indexOp)
	this.ops = append(this.ops, op.doc)
	return this
}

/*
 * Get a string representation of the list of operations in this bulk
 */
func (this *Bulk) String() string {
	ops := ""
	for _, op := range this.ops {
		ops += String(op) + "\n"
	}
	ops = ops[:len(ops)-len("\n")]
	return ops
}

/*
 * Submit a bulk that consists of a list of operations
 * POST /:index/:type/_bulk
 */
func (this *Bulk) Post() {
	log.Println("POST", this.url)
	body := this.String()
	data := bytes.NewReader([]byte(body))
	reader, err := exec("POST", this.url, data)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
