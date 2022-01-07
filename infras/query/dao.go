package query

import (
	"context"
	"github.com/cleango/gallop/infras/errs"
	"gorm.io/gorm"
)

type DAO struct {
	*gorm.DB
}

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func (repo *DAO) ExecTrans(ctx context.Context, fn TransFunc) error {
	if _, ok := FromTrans(ctx); ok {
		return fn(ctx)
	}
	tran := repo.Begin()
	ctx = NewTrans(ctx, tran)
	err := fn(ctx)
	if err != nil {
		_ = tran.Rollback()
		return err
	}
	tran.Commit()
	return nil
}

//GetDB 获取db
func (repo *DAO) GetDB(ctx context.Context) *gorm.DB {
	trans, ok := FromTrans(ctx)
	if ok {
		db, ok := trans.(*gorm.DB)
		if ok {
			return db
		}
	}
	return repo.Begin()
}

//Inserts 批量插入
func (repo *DAO) Save(ctx context.Context, bean interface{}) error {
	err := repo.GetDB(ctx).Create(bean).Error
	return errs.WithDBErr(err)
}

//SaveBatch 批量保存
func (repo *DAO) SaveBatch(ctx context.Context, bean interface{}) error {
	err := repo.GetDB(ctx).CreateInBatches(bean, 100).Error
	return errs.WithDBErr(err)
}
