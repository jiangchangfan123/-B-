package model

import "time"

/*
定义了一个可复用的结构体模板，
配合 GORM 的  AutoMigrate  时，
它会成为数据库表中的对应字段，
并且 GORM 会自动帮你维护创建/更新时间。
*/
type BaseModel struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
