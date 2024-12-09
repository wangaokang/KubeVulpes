package base

import (
	"strconv"
	"time"
)

type Model struct {
	Id              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	GmtCreate       time.Time `gorm:"gmt_create"`
	GmtModified     time.Time `gorm:"gmt_modified"`
	ResourceVersion int64     `gorm:"column:resource_version;not null;default:0" json:"resource_version"`
}

func (m Model) GetSID() string {
	return strconv.FormatInt(m.Id, 10)
}
