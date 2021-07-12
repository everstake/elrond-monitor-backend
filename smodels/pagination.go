package smodels

type Pagination struct {
	Items interface{} `json:"items"`
	Count uint64      `json:"count"`
}
