package elastic

import (
	"fmt"
	"io/ioutil"
	"log"
)

const (
	REFRESH  = "refresh"
	FLUSH    = "flush"
	OPTIMIZE = "optimize"
)

/*
 * A structure for creating shard management operations
 */
type ShardMgmtOp struct {
	url    string
	params map[string]string
}

func newShardMgmtOp(operation string) *ShardMgmtOp {
	return &ShardMgmtOp{url: operation, params: make(map[string]string)}
}

/*
 * Create a refresh API call in order to force recently added document to be visible to search calls
 */
func (this *Elasticsearch) Refresh(index string) *ShardMgmtOp {
	var url string = this.request(index, "", -1, REFRESH)
	return &ShardMgmtOp{url: url}
}

/*
 * Create a flush API call in order to force commit and trauncating the 'translog'
 * See, chapter 11. Inside a shard (Elasticsearch Definitive Guide)
 */
func (this *Elasticsearch) Flush(index string) *ShardMgmtOp {
	var url string = this.request(index, "", -1, FLUSH)
	return &ShardMgmtOp{url: url, params: make(map[string]string)}
}

/*
 * Create an Optimize API call in order to force mering shards into a number of segments
 */
func (this *Elasticsearch) Optimize(index string) *ShardMgmtOp {
	var url string = this.request(index, "", -1, OPTIMIZE)
	return &ShardMgmtOp{url: url, params: make(map[string]string)}
}

/*
 * Add a query parameter to ths Flush API url (e.g. wait_for_ongoing), or Optmize API (e.g. max_num_segment to 1)
 */
func (this *ShardMgmtOp) AddParam(name, value string) *ShardMgmtOp {
	this.params[name] = value
	return this
}

/*
 * Get a string representation of this API url
 */
func (this *ShardMgmtOp) urlString() string {
	return urlString(this.url, this.params)
}

/*
 * Submit a shard managemnt request
 * POST /:index/_refresh
 */
func (this *ShardMgmtOp) Post() {
	url := this.urlString()
	log.Println("POST", url)
	reader, err := exec("POST", this.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
