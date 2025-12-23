package gormUtils

import "gorm.io/gorm"

func Pagination(page int, pageSize int, noLimit bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize <= 0 {
			pageSize = 10
		}

		//如果是不分页，那么直接返回所有
		if noLimit {
			return db
		}

		//这里要减1，因为第1页数据从0开始
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
