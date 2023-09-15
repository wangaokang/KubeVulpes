package types

import "time"

// TimeOption 通用时间规格
type TimeOption struct {
	CreateAt interface{} `json:"create_at,omitempty"`
	UpdateAt interface{} `json:"update_at,omitempty"`
}

const (
	timeLayout = "2006-01-02 15:04:05.999999999"
)

func FormatTime(GmtCreate time.Time, GmtModified time.Time) TimeOption {
	return TimeOption{
		CreateAt: GmtCreate.Format(timeLayout),
		UpdateAt: GmtModified.Format(timeLayout),
	}
}
