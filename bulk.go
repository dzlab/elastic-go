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
func (op *Operation) Add(name string, value interface{}) *Operation {
	op.doc[name] = value
	return op
}

/*
 * Add a field with multiple values to this document
 */
func (op *Operation) AddMultiple(name string, values ...interface{}) *Operation {
	op.doc[name] = values
	return op
}

/*
 * Get a string representation of this operation
 */
func (op *Operation) String() string {
	return String(op.doc)
}

/*
 * Create a new Bulk operations
 */
func (client *Elasticsearch) Bulk(index, docType string) *Bulk {
	url := client.request(index, docType, -1, BULK)
	return &Bulk{url: url, ops: []Dict{}}
}

/*
 * Add an operation to this bulk
 */
func (bulk *Bulk) AddOperation(op *Operation) *Bulk {
	indexOp := Dict{"index": Dict{"_id": op._id}}
	bulk.ops = append(bulk.ops, indexOp)
	bulk.ops = append(bulk.ops, op.doc)
	return bulk
}

/*
 * Get a string representation of the list of operations in this bulk
 */
func (bulk *Bulk) String() string {
	ops := ""
	for _, op := range bulk.ops {
		ops += String(op) + "\n"
	}
	ops = ops[:len(ops)-len("\n")]
	return ops
}

/*
 * Submit a bulk that consists of a list of operations
 * POST /:index/:type/_bulk
 */
func (bulk *Bulk) Post() {
	log.Println("POST", bulk.url)
	body := bulk.String()
	data := bytes.NewReader([]byte(body))
	reader, err := exec("POST", bulk.url, data)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
