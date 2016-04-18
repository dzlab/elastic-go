package elastic

import (
	"encoding/json"
	"log"
)

// Parser an interface for parsing reponses
type Parser interface {
	Parse(data []byte) interface{}
}

// SearchResultParser a parser for search result
type SearchResultParser struct{}

// Parse rerturns a parsed search result structure from the given data
func (parser *SearchResultParser) Parse(data []byte) interface{} {
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

// IndexResultParser a parser for index result
type IndexResultParser struct{}

// Parse returns an index result structure from the given data
func (parser *IndexResultParser) Parse(data []byte) interface{} {
	var result interface{}
	success := Success{}
	if err := json.Unmarshal(data, &success); err == nil {
		log.Println("success", success)
	} else {
		log.Println("Failed to parse response", string(data))
	}
	return result
}

// MappingResultParser a parser for mapping result
type MappingResultParser struct{}

// Parse returns an index result structure from the given data
func (parser *MappingResultParser) Parse(data []byte) interface{} {
	var result interface{}
	/*index := IndexResult{}
	if err := json.Unmarshal(data, &index); err == nil {
		log.Println("index", string(data), index)
	} else {
		log.Println("Failed to parse response", string(data))
	}*/
	log.Println(string(data))
	return result
}

// InsertResultParser a parser for mapping result
type InsertResultParser struct{}

// Parse returns an index result structure from the given data
func (parser *InsertResultParser) Parse(data []byte) interface{} {
	var result interface{}
	insert := InsertResult{}
	if err := json.Unmarshal(data, &insert); err == nil {
		log.Println("insert", string(data), insert)
	} else {
		success := Success{}
		if err := json.Unmarshal(data, &success); err == nil {
			log.Println("success", string(data), success)
			result = success
		} else {
			log.Println("Failed to parse response", string(data))
		}
	}
	return result
}

// AnalyzeResultParser a parser for analyze result
type AnalyzeResultParser struct{}

// Parse returns an analyze result structure from the given data
func (parser *AnalyzeResultParser) Parse(data []byte) interface{} {
	var result interface{}
	analyze := AnalyzeResult{}
	if err := json.Unmarshal(data, &analyze); err == nil {
		log.Println("analyze", string(data), analyze)
	} else {
		success := Success{}
		if err := json.Unmarshal(data, &success); err == nil {
			log.Println("success", string(data), success)
			result = success
		} else {
			log.Println("Failed to parse response", string(data))
		}
	}
	return result
}

// BulkResultParser a parser for analyze result
type BulkResultParser struct{}

// Parse returns an analyze result structure from the given data
func (parser *BulkResultParser) Parse(data []byte) interface{} {
	var result interface{}
	bulk := BulkResult{}
	if err := json.Unmarshal(data, &bulk); err == nil {
		log.Println("bulk", string(data), bulk)
	} else {
		success := Success{}
		if err := json.Unmarshal(data, &success); err == nil {
			log.Println("success", string(data), success)
			result = success
		} else {
			log.Println("Failed to parse response", string(data))
		}
	}
	return result
}
