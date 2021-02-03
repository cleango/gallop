package page

import (
	"github.com/cleango/gallop/infras/errs"
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
	GetPage() Result
}

//Params 分页参数
type Params struct {
	PageIndex int    `json:"page_index" form:"page_index"`
	PageSize  int    `json:"page_size" form:"page_size"`
	Asc       string `json:"asc" form:"asc"`
	Desc      string `json:"desc" form:"desc"`
}

//GetPageIndex 获取当前页码
func (p *Params) GetPageIndex() int {
	if p.PageIndex == 0 {
		return 1
	}
	return p.PageIndex
}

//GetPageSize 获取分页大小
func (p *Params) GetPageSize() int {
	if p.PageSize == 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

//GetPage 获取分页
func (p *Params) GetPage() Result {
	return Result{
		PageIndex: p.GetPageIndex(),
		PageSize:  p.GetPageSize(),
	}
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
	PageIndex int   `json:"page_index"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
}

//Execute 分页查询
func Execute(db *xorm.Session, query Query, bean interface{}) (res Result, err error) {
	res = Result{}
	defer func() {
		err = errs.WithDBErr(err)
	}()
	sql, args, err := builder.ToSQL(query.Build())
	if err != nil {
		return res, errs.WithDBErr(err)
	}
	res.PageIndex = query.GetPageIndex()
	res.PageSize = query.GetPageSize()

	if sql != "" {
		db = db.Where(sql, args...)
	}
	if res.PageSize == -1 {
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
			return
		}
		return
	}
	db = db.Limit(query.GetPageSize(), (res.PageIndex-1)*res.PageSize)
	asc := query.GetAsc()
	if len(asc) > 0 {
		db = db.Asc(asc...)
	}
	desc := query.GetDesc()
	if len(desc) > 0 {
		db = db.Desc(desc...)
	}
	res.Total, err = db.FindAndCount(bean)
	if err != nil {
		return res, errs.WithDBErr(err)
	}
	return
}

//ExecAll 不分页查询多条
func ExecAll(db *xorm.Session, query builder.Cond, bean interface{}) (err error) {
	sql, args, err := builder.ToSQL(query)
	if err != nil {
		return err
	}
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err = db.Find(bean)
	if err != nil {
		return errs.WithDBErr(err)
	}
	return
}

//ExecOne 动态查询获取一条数据
func ExecOne(db *xorm.Session, query builder.Cond, bean interface{}) (has bool, err error) {
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
