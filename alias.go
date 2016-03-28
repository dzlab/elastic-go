package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	ALIASES = "_aliases"
	ACTIONS = "actions"
)

type Alias struct {
	url  string
	dict Dict
}

/*
 * Return a JSOn representation of the body of this Alias
 */
func (this *Alias) String() string {
	result, err := json.Marshal(this.dict)
	if err != nil {
		log.Println(err)
	}
	return string(result)
}

/*
 * Create an alias
 */
func newAlias() *Alias {
	return &Alias{url: "", dict: make(Dict)}
}

func (this *Elasticsearch) Alias() *Alias {
	url := fmt.Sprintf("http://%s/%s", this.Addr, ALIASES)
	return &Alias{url: url, dict: make(Dict)}
}

/*
 * Add an Alias operation (e.g. remove index's alias)
 */
func (this *Alias) AddAction(operation, index, alias string) *Alias {
	if this.dict[ACTIONS] == nil {
		this.dict[ACTIONS] = []Dict{}
	}
	action := make(Dict)
	action[operation] = Dict{"index": index, "alias": alias}
	this.dict[ACTIONS] = append(this.dict[ACTIONS].([]Dict), action)
	return this
}

/*
 * Submit an Aliases POST operation
 * POST /:index
 */
func (this *Alias) Post() {
	log.Println("POST", this.url)
	body := String(this.dict)
	reader, err := exec("POST", this.url, bytes.NewReader([]byte(body)))
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}

}
