# esSearchBuilder
a simple esSearchBuilder for php

es index mapping
```php
{
	"index": "goods",
	"body": {
		"settings": {
			"number_of_shards": 3,
			"number_of_replicas": 2,
			"analysis": {
				"analyzer": {
					"default": {
						"tokenizer": "ik_max_word"
					},
					"pinyin_analyzer": {
						"tokenizer": "my_pinyin"
					}
				},
				"tokenizer": {
					"my_pinyin": {
						"type": "pinyin",
						"keep_first_letter": true,
						"keep_separate_first_letter": false,
						"keep_full_pinyin": true,
						"keep_original": true,
						"limit_first_letter_length": 16,
						"lowercase": true,
						"remove_duplicated_term": true
					}
				}
			}
		},
		"mappings": {
			"_source": {
				"enabled": true
			},
			"properties": {
				"goods_id": {
					"type": "integer"
				},
				"goods_name": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word",
					"fields": {
						"pinyin": {
							"type": "text",
							"term_vector": "with_positions_offsets",
							"analyzer": "pinyin_analyzer"
						}
					}
				},
				"goods_price": {
					"type": "double"
				},
				"goods_standards": {
					"type": "text"
				},
				"goods_indications": {
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"goods_image": {
					"type": "text"
				},
				"store_name": {
					"type": "text"
				},
				"store_self_pickup": {
					"type": "integer"
				},
				"store_no_rest": {
					"type": "integer"
				},
				"store_free_freight": {
					"type": "integer"
				},
				"store_id": {
					"type": "integer"
				},
				"store_delivery": {
					"type": "integer"
				},
				"gc_id_3": {
					"type": "integer"
				},
				"location": {
					"type": "geo_point"
				},
				"city_id": {
					"type": "integer"
				},
				"hospital_id": {
					"type": "integer"
				}
			}
		}
	}
}
```

demo
```php
// 关键字
        $keywords = $this->request->input('keywords', '');
        // 页码
        $page = $this->request->input('page');
        // 每页显示条数
        $pageSize = $this->request->input('page_size');

        $builder = (new EsSearchBuilder('goods'))->paginate((int) $pageSize, (int) $page);

        // 假如有关键词查询，拼装查询条件
        if ($keywords) {
            $keywordsArr = array_filter(explode(' ', $keywords));
            $multiMatch = ['goods_name^3', 'goods_name.pinyin^2', 'goods_indications'];
            $builder->keywords($keywordsArr, $multiMatch)->sort('_score', ['order' => 'desc']);
        }

        // 拼装查询条件
        $regx = [0, 1];

        if (in_array($this->request->input('store_self_pickup', ''), $regx)) {
            $builder->matchValue('store_self_pickup', $this->request->input('store_self_pickup'));
        }

        if (in_array($this->request->input('store_no_rest', ''), $regx)) {
            $builder->matchValue('store_no_rest', $this->request->input('store_no_rest'));
        }

        if (in_array($this->request->input('store_free_freight', ''), $regx)) {
            $builder->matchValue('store_free_freight', $this->request->input('store_free_freight'));
        }

        if (in_array($this->request->input('store_delivery', ''), $regx)) {
            $builder->matchValue('store_delivery', $this->request->input('store_delivery'));
        }

        if ($this->request->input('gc_id_3')) {
            $builder->matchValue('gc_id_3', $this->request->input('gc_id_3'));
        }

        if ($this->request->input('hospital_id')) {
            $builder->matchValue('hospital_id', $this->request->input('hospital_id'));
        }

        // 有经纬度，则搜索经纬度所在的城市的商品，然后按照距离排序
        $location = $this->request->input('location', '');
        if ($location) {
            $location = explode(',', $location);
            $cityCode = $this->request->input('city_code');
            // 这是小程序获取地理位置时，cityCode是156400100这个格式，所以稍微处理一下，具体情况具体分析
            $cityCode = str_replace('156', '', $cityCode);
            $builder->matchValue('city_id', $cityCode)->sort('_geo_distance', [
                'unit' => 'km',
                'location' => [
                    'lon' => $location[0],
                    'lat' => $location[1],
                ],
                'order' => 'asc',
            ]);

            // 如果有距离范围要求，则启用下面代码
            if ($this->request->input('search_km')) {
                $builder->filter([
                    'geo_distance' => [
                        'distance' => $this->request->input('search_km', 1) . 'km',
                        'location' => [
                            'lon' => $location[0],
                            'lat' => $location[1],
                        ],
                    ],
                ]);
            }
        }

        $priceRange = [];
        if ($this->request->input('price_max')) {
            $priceRange['range']['goods_price']['lte'] = $this->request->input('price_max');
        }

        if ($this->request->input('price_min')) {
            $priceRange['range']['goods_price']['gte'] = $this->request->input('price_min');
        }

        if ($priceRange) {
            $builder->filter($priceRange);
        }

        $query = $builder->getParams();

        return $this->esClient->search($query);
```

# ChangeLog
- 2022/09/28 新增go版本，详情请看GoVersion，使用的库为github.com/elastic/go-elasticsearch