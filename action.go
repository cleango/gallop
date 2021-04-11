package gallop

import "context"

type IAction interface {
	Exec()
}

type IActionClose interface {
}
type CloseContext context.Context
type IClose interface {
	Shutdown(CloseContext)
}
