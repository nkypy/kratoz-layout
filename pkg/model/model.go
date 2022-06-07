package model

import "time"

// Model a basic GoLang struct which includes the following fields: SID, UUID, CreateTime
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      Model
//    }
type Model struct {
	ID         uint64    `gorm:"column:Sid"`  // 自增序号
	UUID       string    `gorm:"column:UUID"` // 唯一识别码
	CreateTime time.Time `gorm:"autoUpdateTime:nano"`
}
