package repository

import (
	"github.com/irvankadhafi/erajaya-product-service/utils"
	"gorm.io/gorm"
)

// scopeByPageAndLimit return a gorm offset clause by page and limit
func scopeByPageAndLimit(page, limit int) func(d *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Offset(utils.Offset(page, limit)).Limit(limit) }
}
