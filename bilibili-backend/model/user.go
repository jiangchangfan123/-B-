package model

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(32);not null;uniqueIndex" json:"username"`
	Nickname string `gorm:"type:varchar(32);default:''" json:"nickname"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(128);uniqueIndex" json:"email"`
	Avatar   string `gorm:"type:varchar(255);default:''" json:"avatar"`
	Sign     string `gorm:"type:varchar(200);default:''" json:"sign"`
	Role     int    `gorm:"type:tinyint;default:1" json:"role"`
	Coins    int    `gorm:"type:int;default:0" json:"coins"`
	Status   int    `gorm:"type:tinyint;default:1" json:"status"`
}
