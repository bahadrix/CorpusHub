package index

import "time"

type SearchResult struct {
	TotalHits uint64 `json:"total_hits"`
	MaxScore float64 `json:"max_score"`
	Took time.Duration `json:"took"`
	Hits []*SearchHit `json:"hits"`
}

type SearchHit struct {
	ID string `json:"id"`
	Score float64 `json:"score"`
	Fragments map[string][]string `json:"fragments,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}
