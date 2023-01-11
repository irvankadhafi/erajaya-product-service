package cache

import "time"

const (
	defaultTTL    = 10 * time.Second
	defaultNilTTL = 5 * time.Minute
)
