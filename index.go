package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	SHARDS_NB   = "number_of_shards"
	REPLICAS_NB = "number_of_replicas"
	ANALYSIS    = "analysis"
	SETTINGS    = "settings"
)

type Index struct {
	url  string
	dict Dict
	//settings *Settings
	//mappings map[string]string
}

/*
 * Return a JSOn representation of the body of this Index
 */
func (this *Index) String() string {
	result, err := json.Marshal(this.dict)
	if err != nil {
		log.Println(err)
	}
	return string(result)
}

func (this *Elasticsearch) Index(index string) *Index {
	url := fmt.Sprintf("http://%s/%s", this.Addr, index)
	return &Index{url: url, dict: make(Dict)}
}

/*
 * add a setting parameter
 */
func (this *Index) Settings(settings Dict) {
	this.dict[SETTINGS] = settings
}

/*
 * Create new index settings
 */
func newIndex() *Index {
	return &Index{dict: make(Dict)}
}

/*
 * Add a key-value settings
 */
func (this *Index) AddSetting(name string, value int) *Index {
	if this.dict[SETTINGS] == nil {
		this.dict[SETTINGS] = make(Dict)
	}
	this.dict[SETTINGS].(Dict)[name] = value
	return this
}

/*
 * Set the number of shards
 */
func (this *Index) SetShardsNb(number int) *Index {
	this.AddSetting(SHARDS_NB, number)
	return this
}

/*
 * Set the number of shards
 */
func (this *Index) SetReplicasNb(number int) *Index {
	this.AddSetting(REPLICAS_NB, number)
	return this
}

/*
 * Analyzer/Filter
 */
type Analyzer struct {
	name string
	kv   map[string]Dict
}

/*
 * Create a new analyzer
 */
func NewAnalyzer(name string) *Analyzer {
	return &Analyzer{name: name, kv: make(map[string]Dict)}
}

/*
 * Return a JSON string representation of this analyzer
 */
func (this *Analyzer) String() string {
	dict := make(Dict)
	dict[this.name] = this.kv
	return String(dict)
}

/*
 * Add an anlyzer to the index settings
 */
func (this *Index) AddAnalyzer(analyzer *Analyzer) *Index {
	// if no "settings" create one
	if this.dict[SETTINGS] == nil {
		this.dict[SETTINGS] = make(Dict)
	}
	// if no "settings.analysis" create one
	if this.dict[SETTINGS].(Dict)[ANALYSIS] == nil {
		this.dict[SETTINGS].(Dict)[ANALYSIS] = make(Dict) //map[string]*Analyzer)
	}
	// insert the analyser ('name' and 'kv' attributes are taken separately)
	settings := this.dict[SETTINGS].(Dict)
	analysis := settings[ANALYSIS].(Dict) //map[string]*Analyzer)
	analysis[analyzer.name] = analyzer.kv
	this.dict[SETTINGS].(Dict)[ANALYSIS] = analysis
	return this
}

/*
 * add an attribute to analyzer definition
 */
func (this *Analyzer) Add1(key1, key2 string, value interface{}) *Analyzer {
	if len(this.kv[key1]) == 0 {
		this.kv[key1] = make(Dict)
	}
	this.kv[key1][key2] = value
	return this
}

/*
 * add a dictionary of attributes to analyzer definition
 */
func (this *Analyzer) Add2(name string, value Dict) *Analyzer {
	if len(this.kv[name]) == 0 {
		this.kv[name] = make(Dict)
	}
	for k, v := range value {
		this.kv[name][k] = v
	}
	return this
}

/*
 * Create an index
 * PUT /:index
 */
func (this *Index) Put() {
	log.Println("PUT", this.url)
	body := String(this.dict)
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
 * delete an index
 * DELETE /:index
 */
func (this *Index) Delete() {
	reader, err := exec("DELETE", this.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
