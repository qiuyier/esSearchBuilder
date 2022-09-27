package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	cert, _ := os.ReadFile("path/to/http_ca.crt")
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://127.0.0.1:9200",
		},
		Username: "*****",
		Password: "*******",
		CACert:   cert,
	}

	es, _ := elasticsearch.NewClient(cfg)

	//fmt.Println(elasticsearch.Version)
	//fmt.Println(es.Info())
	var searchBuilder = &SearchBuilder{}
	searchBuilder.matchValue("store_self_pickup", 0)
	searchBuilder.matchValue("store_no_rest", 1)
	keywords := "奇"
	keywordsArr := strings.Split(keywords, " ")
	searchBuilder.keywords(keywordsArr, []string{"goods_name^3", "goods_name.pinyin^2", "goods_indications"})
	searchBuilder.sort("_score", map[string]interface{}{"order": "desc"}).paginate(10, 1)
	body, err := searchBuilder.getParamJson()
	if err != nil {
		fmt.Println("getParamJson err: ", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("goods"),
		es.Search.WithBody(&body),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		fmt.Println("es err", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("close: ", err)
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			fmt.Println("Error parsing the response body: ", err)
		} else {
			// Print the response status and error information.
			fmt.Printf("[%s] %s: %s\n",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	r := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("Error parsing the response body: ", err)
	}

	fmt.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		fmt.Printf(" * ID=%s, %s\n", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	fmt.Println(strings.Repeat("=", 37))
}
