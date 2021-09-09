package models

// Price is structure for Price objects
type Price struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Bid  uint64 `json:"bid"`
	Ask  uint64 `json:"ask"`
}
