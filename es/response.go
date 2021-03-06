package es

import "github.com/dynport/dgtk/es/aggregations"

func NewResponse(raw []byte) *Response {
	return &Response{
		Raw:          raw,
		Aggregations: aggregations.Aggregations{},
	}
}

type ResponseRaw struct {
	*Response
	Hits *HitsRaw `json:"hits"`
}

type Response struct {
	Took         int                       `json:"took"`
	TimedOut     bool                      `json:"timed_out"`
	Facets       ResponseFacets            `json:"facets"`
	Hits         Hits                      `json:"hits"`
	Shards       *ShardsResponse           `json:"_shards,omitempty"`
	Raw          []byte                    `json:"-"`
	Aggregations aggregations.Aggregations `json:"aggregations,omitempty"`
}

func (s *Response) SetRaw(b []byte) {
	s.Raw = b
}

func (r *Response) ShardsResponse() *ShardsResponse {
	return r.Shards
}

type ShardsResponse struct {
	Total      int        `json:"total"`
	Successful int        `json:"successful"`
	Failed     int        `json:"failed"`
	Failures   []*Failure `json:"failures"`
}

type Failure struct {
	Index  string `json:"index"`
	Shard  int    `json:"shard"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

type ResponseFacets map[string]*ResponseFacet

type ResponseFacet struct {
	Type    string       `json:"_type"`
	Missing int          `json:"missing"`
	Total   int          `json:"total"`
	Other   int          `json:"other"`
	Terms   []*FacetTerm `json:"terms,omitempty"`
	Entries []*Entry     `json:"entries,omitempty"`
}

type FacetTerm struct {
	Term  interface{} `json:"term"`
	Count int         `json:"count"`
}
