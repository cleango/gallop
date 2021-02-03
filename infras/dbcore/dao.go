package dbcore

import (
	"context"
	"github.com/cleango/gallop/infras/errs"
	"reflect"
	"xorm.io/xorm"
)

type DAO struct {
	*xorm.Engine
}

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func (repo *DAO) ExecTrans(ctx context.Context, fn TransFunc) error {
	if _, ok := FromTrans(ctx); ok {
		return fn(ctx)
	}
	tran := repo.Engine.NewSession()

	defer tran.Close()
	if err := tran.Begin(); err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tran.Rollback()
			panic(r)
		}
	}()

	ctx = NewTrans(ctx, tran)
	err := fn(ctx)
	if err != nil {
		_ = tran.Rollback()
		return err
	}
	return tran.Commit()
}

//GetDB 获取db
func (repo *DAO) GetDB(ctx context.Context) *xorm.Session {
	trans, ok := FromTrans(ctx)
	if ok {
		db, ok := trans.(*xorm.Session)
		if ok {
			return db
		}
	}
	return repo.Engine.NewSession()
}

//Inserts 批量插入
func (repo *DAO) Save(ctx context.Context, bean interface{}) error {
	_, err := repo.GetDB(ctx).Insert(bean)
	return errs.WithDBErr(err)
}

//SaveBatch 批量保存
func (repo *DAO) SaveBatch(ctx context.Context, bean interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(bean))
	if sliceValue.Kind() != reflect.Slice {
		return errs.WithDBErrf("needs a pointer to a slice")
	}

	if sliceValue.Len() <= 0 {
		return errs.WithDBErrf("could not insert a empty slice")
	}
	var (
		size = sliceValue.Len()
	)
	return repo.ExecTrans(ctx, func(ctx context.Context) error {
		var temp []interface{}
		for i := 0; i < size; i++ {
			v := sliceValue.Index(i)
			var vv reflect.Value
			switch v.Kind() {
			case reflect.Interface:
				vv = reflect.Indirect(v.Elem())
			default:
				vv = reflect.Indirect(v)
			}
			temp = append(temp, vv.Interface())
			if len(temp) == 1000 {
				err := repo.Save(ctx, &temp)
				if err != nil {
					return err
				}
				temp = make([]interface{}, 0)
			}
		}
		if len(temp) > 0 {
			err := repo.Save(ctx, &temp)
			if err != nil {
				return errs.WithDBErr(err)
			}
		}
		return nil
	})

}

//Update 报错
func (repo *DAO) UpdateByID(ctx context.Context, bean interface{}, id interface{}, cols ...string) error {
	db := repo.GetDB(ctx).ID(id)
	if len(cols) > 0 {
		db = db.Cols(cols...)
	}
	_, err := db.Update(bean)
	return errs.WithDBErr(err)
}
