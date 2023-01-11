package repository

import (
	"encoding/json"
	"github.com/irvankadhafi/erajaya-product-service/cache"
	"github.com/irvankadhafi/erajaya-product-service/internal/model"
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// scopeByPageAndLimit is a helper function to apply pagination on gorm query.
// it takes in 2 input as page and limit and returns a scope function
// that can be passed to gorm's db.Scopes method
// It is reusable to apply pagination on any query where it is needed
func scopeByPageAndLimit(page, limit int) func(d *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Offset(utils.Offset(page, limit)).Limit(limit) }
}

// scopeMatchTSQuery is a helper function to apply plainto_tsquery on gorm query.
// it takes in a query string as input and returns a scope function
// that can be passed to gorm's db.Scopes method
// It is reusable to apply plainto_tsquery on any query where it is needed
func scopeMatchTSQuery(query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name @@ plainto_tsquery(?)", query)
	}
}

// storeNilCacheByKey store nil value to cache
func storeNilCacheByKey(c cache.Cache, key string) {
	if err := c.StoreNil(key); err != nil {
		logrus.Error(err)
	}
}

// storeNilMultiValueByKey store nil multi value to cache
func storeNilMultiValueByKey(c cache.Cache, key string) {
	cacheItem := cache.NewItemWithCustomTTL(
		key,
		utils.Dump(model.MultiCacheValue{}),
		15*time.Second,
	)
	if err := c.Store(cacheItem); err != nil {
		logrus.Error(err)
	}
}

// findFromCacheByKey :nodoc:
func findFromCacheByKey[T any](c cache.Cache, key string) (reply T, err error) {
	var rep any
	rep, err = c.Get(key)
	if err != nil || rep == nil {
		return
	}

	bt, _ := rep.([]byte)
	if bt == nil {
		return
	}

	if err = json.Unmarshal(bt, &reply); err != nil {
		return
	}

	return
}
