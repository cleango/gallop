package query

import (
	"fmt"
	"github.com/cleango/gallop/infras/errs"
	"gorm.io/gorm"
	"math"
)

//Info 分页查询条件
type Query interface {
	GetPageIndex() int
	GetPageSize() int
	Build(db *gorm.DB) *gorm.DB
	GetOrder() string
}

//Params 分页参数
type Params struct {
	Page  int    `json:"page" form:"page"`
	Limit int    `json:"limit" form:"limit"`
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}

//GetPageIndex 获取当前页码
func (p *Params) GetPageIndex() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

//GetPageSize 获取分页大小
func (p *Params) GetPageSize() int {
	if p.Limit == 0 {
		return 20
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

//GetOrder 获取排序
func (p *Params) GetOrder() string {
	if len(p.Sort) > 0 && len(p.Order) > 0 {
		return fmt.Sprintf("%s %s", p.Sort, p.Order)
	}
	return ""
}

//Data 分页结果
type Result struct {
	Records interface{} `json:"records"`
	Current int         `json:"current"`
	Pages   int         `json:"pages"`
	Size    int         `json:"size"`
	Total   int         `json:"total"`
}

func (r *Result) WithRecords(data interface{}) *Result {
	r.Records = data
	return r
}
func NewResult(page, limit, total int) *Result {
	pages := int(math.Ceil(float64(total) / float64(limit)))
	return &Result{
		Records: nil,
		Current: page,
		Pages:   pages,
		Size:    limit,
		Total:   total,
	}
}

//Execute 分页查询
func Page(db *gorm.DB, query Query, bean interface{}) (*Result, error) {
	db = query.Build(db)

	pageIndex := query.GetPageIndex()
	pageSize := query.GetPageSize()
	var total int64

	res := db.Count(&total)
	if res.Error != nil {
		return nil, errs.WithDBErr(res.Error)
	}

	db = db.Offset((pageIndex - 1) * pageSize).Limit(query.GetPageSize())
	order := query.GetOrder()
	if len(order) > 0 {
		db = db.Order(order)
	}

	result := db.Find(bean)
	if result.Error != nil {
		return nil, errs.WithDBErr(result.Error)
	}

	return NewResult(pageIndex, pageSize, int(total)).WithRecords(bean), nil
}

//Exec 不分页查询多条
func Exec(db *gorm.DB, query Query, bean interface{}) (err error) {
	db = query.Build(db)
	order := query.GetOrder()
	if len(order) > 0 {
		db = db.Order(order)
	}
	res := db.Find(bean)
	if res.Error != nil {
		return errs.WithDBErr(res.Error)
	}
	return
}
