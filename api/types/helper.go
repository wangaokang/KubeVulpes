package types

// TimeOption 通用时间规格
type TimeOption struct {
	CreateAt interface{} `json:"create_at,omitempty"`
	UpdateAt interface{} `json:"update_at,omitempty"`
}
