package repository

import (
	"gorm.io/gorm"
	"strings"
)

func orderByCreatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}

func orderByCreatedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at ASC")
}

func orderByUpdatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at DESC")
}

func orderByUpdatedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at ASC")
}

func withSize(size int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(size))
	}
}

func withOffset(offset int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset))
	}
}

func withNameLike(query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(query)+"%")
	}
}
