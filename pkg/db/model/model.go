package model

import "time"

type KubeVulpes struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	IsDeleted int       `gorm:"column:is_deleted" json:"is_deleted"`
}

type User struct {
	KubeVulpes
	Name        string `gorm:"index:idx_name,unique" json:"name"`
	Password    string `gorm:"type:varchar(256)" json:"password"`
	Role        string `gorm:"column:role;not null" json:"role"`
	Status      int8   `gorm:"column:status;not null" json:"status"`
	Email       string `gorm:"type:varchar(128),unique" json:"email"`
	Description string `gorm:"type:text" json:"description"`
}

func (user *User) TableName() string {
	return "users"
}
