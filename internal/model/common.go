package model

// MultiCacheValue store multiple values for bucket caching
type MultiCacheValue struct {
	IDs   []int64 `json:"ids"`
	Count int64   `json:"count"`
}
