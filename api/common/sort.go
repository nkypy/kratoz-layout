package common

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"switch_data_center_go/pkg/util"
)

func (x *SortRequest) Values() string {
	if x == nil {
		return ""
	}

	if !x.isFieldValid() {
		return ""
	}

	if util.Empty(x.Field) || util.Empty(x.Order) {
		return ""
	}

	return fmt.Sprintf("%s %s", x.Field, x.expression())
}

func (x *SortRequest) expression() string {
	if x == nil {
		return ""
	}

	switch strings.ToLower(x.Order) {
	case "asc", "ascend":
		return "ASC"
	case "desc", "descend":
		return "DESC"
	default:
		return ""
	}
}

// 返回 Field 值是否有效(防止 sql 注入等问题)
func (x *SortRequest) isFieldValid() bool {
	if x == nil {
		return false
	}

	if x.Field == "" {
		return false
	}
	return !strings.ContainsAny(x.Field, " &`'\"?<>:;()[]|@=-#~!$%^*=")
}

// Sort specify order when retrieve records from database
func (x *SortRequest) Sort(db *gorm.DB, defaultSortValues ...string) {
	if x == nil {
		return
	}

	sortValue := x.Values()
	if !util.Empty(sortValue) {
		db.Order(sortValue)
		return
	}

	// 默认排序
	if len(defaultSortValues) > 0 {
		for _, value := range defaultSortValues {
			db.Order(value)
		}
	}
}
