package elastic

import (
	"encoding/json"
	"log"
)

/*
 * An interface for parsing reponses
 */
type Parser interface {
	Parse(data []byte) interface{}
}

/*
 * A parser for search result
 */
type SearchResultParser struct{}

/*
 * Parse the given data into a search result structure
 */
func (this *SearchResultParser) Parse(data []byte) interface{} {
	var result interface{}
	search := SearchResult{}
	if err := json.Unmarshal(data, &search); err == nil {
		log.Println("search", string(data), search)
		result = search
	} else {
		success := Success{}
		if err := json.Unmarshal(data, &success); err == nil {
			log.Println("success", string(data), success)
			result = success
		} else {
			failure := Failure{}
			if err := json.Unmarshal(data, &failure); err == nil {
				log.Println("failed", string(data), failure)
				result = failure
			} else {
				log.Println("Failed to parse response", string(data))
			}
		}
	}
	return result
}
