package models

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func Paginate(req *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := req.URL.Query()
		limit := 10
		offset := 0
		if strLimit := query.Get("limit"); strLimit != "" {
			if lmt, err := strconv.Atoi(strLimit); err == nil {
				limit = lmt
			}
		}
		if strOffset := query.Get("offset"); strOffset != "" {
			if ofst, err := strconv.Atoi(strOffset); err == nil {
				offset = ofst
			}
		}
		return db.Offset(offset).Limit(limit)
	}
}
