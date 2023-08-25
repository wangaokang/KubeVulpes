package meta

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	offset = 1  //默认取第一页
	limit  = 10 //默认取10条记录
)

type ListSelector struct {
	Page  int // 页数
	Limit int // 每页数量
}

func ParseListSelector(c *gin.Context) *ListSelector {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = offset
	}
	pageSize, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		pageSize = limit
	}

	ls := &ListSelector{
		Page:  page,
		Limit: pageSize,
	}
	return ls
}
