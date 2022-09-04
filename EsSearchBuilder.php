<?php

declare(strict_types=1);
/**
 * User: qiuyier
 * Date: 2022/9/5
 * TIME: 00:36
 */
class EsSearchBuilder
{
    protected string $index;

    protected array $params;

    public function __construct(string $index, string $doc = '')
    {
        $this->index = $index;
        $this->params = [
            'index' => $this->index,
            'body' => [
                'query' => [
                    'bool' => [
                        'filter' => [],
                        'must' => [],
                    ],
                ],
            ],
        ];
        if ($doc) {
            $this->params['type'] = $doc;
        }
    }

    /**
     * @return $this
     */
    public function paginate(int $size, int $page): static
    {
        $this->params['body']['from'] = ($page - 1) * $size;
        $this->params['body']['size'] = $size;

        return $this;
    }

    /**
     * @param $key
     * @param $value
     * @return $this
     */
    public function matchValue($key, $value): static
    {
        $this->params['body']['query']['bool']['filter'][] = ['term' => [$key => $value]];

        return $this;
    }

    /**
     * @param $filter
     * @return $this
     */
    public function filter($filter): static
    {
        $this->params['body']['query']['bool']['filter'][] = $filter;

        return $this;
    }

    /**
     * @param $keywords
     * @param $matchField ['title^3','long_title^2','category^2','description','skus_title','skus_description','properties_value',]
     * @return $this
     */
    public function keywords($keywords, $matchField): static
    {
        $keywords = is_array($keywords) ? $keywords : [$keywords];

        foreach ($keywords as $keyword) {
            $this->params['body']['query']['bool']['must'][] = [
                'multi_match' => [
                    'query' => $keyword,
                    'fields' => $matchField,
                ],
            ];
        }

        return $this;
    }

    /**
     * @param $field
     * @param $direction
     * @return $this
     */
    public function sort($field, $direction): static
    {
        if (! isset($this->params['body']['sort'])) {
            $this->params['body']['sort'] = [];
        }
        $this->params['body']['sort'][] = [$field => $direction];

        return $this;
    }

    public function getParams(): array
    {
        return $this->params;
    }
}