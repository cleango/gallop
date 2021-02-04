package page

import (
	"github.com/cleango/gallop/infras/errs"
	"math"
	"strings"

	"xorm.io/builder"
	"xorm.io/xorm"
)

//Info 分页查询条件
type Query interface {
	GetPageIndex() int
	GetPageSize() int
	Build() builder.Cond
	GetAsc() []string
	GetDesc() []string
}

//Params 分页参数
type Params struct {
	Page  int    `json:"page" form:"page"`
	Limit int    `json:"limit" form:"limit"`
	Asc   string `json:"asc" form:"asc"`
	Desc  string `json:"desc" form:"desc"`
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

//GetAsc 获取正序
func (p *Params) GetAsc() []string {
	s := strings.TrimSpace(p.Asc)
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

//GetDesc 获取倒叙
func (p *Params) GetDesc() []string {
	s := strings.TrimSpace(p.Desc)
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
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
func Execute(db *xorm.Session, query Query, bean interface{}) (*Result, error) {
	sql, args, err := builder.ToSQL(query.Build())
	if err != nil {
		return nil, errs.WithDBErr(err)
	}
	pageIndex := query.GetPageIndex()
	pageSize := query.GetPageSize()

	if sql != "" {
		db = db.Where(sql, args...)
	}
	db = db.Limit(query.GetPageSize(), (pageIndex-1)*pageSize)
	asc := query.GetAsc()
	if len(asc) > 0 {
		db = db.Asc(asc...)
	}
	desc := query.GetDesc()
	if len(desc) > 0 {
		db = db.Desc(desc...)
	}
	var total int64
	total, err = db.FindAndCount(bean)
	if err != nil {
		return nil, errs.WithDBErr(err)
	}
	return NewResult(pageIndex, pageSize, int(total)).WithRecords(bean), nil
}

//ExecAll 不分页查询多条
func ExecAll(db *xorm.Session, query Query, bean interface{}) (err error) {
	sql, args, err := builder.ToSQL(query)
	if err != nil {
		return err
	}
	if sql != "" {
		db = db.Where(sql, args...)
	}
	asc := query.GetAsc()
	if len(asc) > 0 {
		db = db.Asc(asc...)
	}
	desc := query.GetDesc()
	if len(desc) > 0 {
		db = db.Desc(desc...)
	}
	err = db.Find(bean)
	if err != nil {
		return errs.WithDBErr(err)
	}
	return
}

//ExecOne 动态查询获取一条数据
func ExecOne(db *xorm.Session, query Query, bean interface{}) (has bool, err error) {

	sql, args, err := builder.ToSQL(query)
	if err != nil {
		return false, errs.WithDBErr(err)
	}
	if sql != "" {
		db = db.Where(sql, args...)
	}
	has, err = db.Get(bean)
	if err != nil {
		return false, errs.WithDBErr(err)
	}
	return has, nil
}
