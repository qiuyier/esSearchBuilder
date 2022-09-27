package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SearchBuilder struct {
	SearchCollection
}

type SearchCollection struct {
	Query struct {
		Bool struct {
			Filter []interface{} `json:"filter"`
			Must   []interface{} `json:"must"`
		} `json:"bool"`
	} `json:"query"`
	From int           `json:"from"`
	Size int           `json:"size"`
	Sort []interface{} `json:"sort"`
}

func (search *SearchCollection) paginate(size, page int) *SearchCollection {
	search.From = (page - 1) * size
	search.Size = size
	return search
}

func (search *SearchCollection) matchValue(key string, value interface{}) *SearchCollection {
	query := map[string]interface{}{
		"term": map[string]interface{}{
			key: value,
		},
	}
	//filter = append(filter, query)
	search.Query.Bool.Filter = append(search.Query.Bool.Filter, query)
	return search
}

func (search *SearchCollection) keywords(keywords, matchFields []string) *SearchCollection {
	for _, keyword := range keywords {
		query := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": matchFields,
			},
		}
		//must = append(must, query)
		search.Query.Bool.Must = append(search.Query.Bool.Must, query)
	}

	return search
}

func (search *SearchCollection) sort(field string, direction interface{}) *SearchCollection {
	query := map[string]interface{}{
		field: direction,
	}
	//sort = append(sort, query)
	search.Sort = append(search.Sort, query)
	return search
}

func (search *SearchCollection) getParamJson() (bytes.Buffer, error) {
	p, err := json.Marshal(search)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(p))
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(search); err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}
