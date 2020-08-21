package gallop

import (
	"fmt"
)

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
func ThrowWithErr(err, base error) {
	if err != nil {
		panic(fmt.Errorf("%w %s", base, err.Error()))
	}
}

