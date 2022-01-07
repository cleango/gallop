package query

import "context"

// 定义全局上下文中的键
type (
	transCtx     struct{}
	transLockCtx struct{}
)

// NewTrans 创建事务的上下文
func NewTrans(ctx context.Context, trans interface{}) context.Context {
	return context.WithValue(ctx, transCtx{}, trans)
}

// FromTrans 从上下文中获取事务
func FromTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(transCtx{})
	return v, v != nil
}

// NewTransLock 创建事务锁的上下文
func NewTransLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, transLockCtx{}, struct{}{})
}

// FromTransLock 从上下文中获取事务锁
func FromTransLock(ctx context.Context) bool {
	v := ctx.Value(transLockCtx{})
	return v != nil
}
