package util

import (
	"fmt"
)

const (
	KESuccess int32 = 99999 - iota
)

type TaskDevError struct {
	Code int32
}

var err map[int32]string

func (e TaskDevError) Error() string {
	if msg, ok := err[e.Code]; ok {
		return fmt.Sprintf("[ERR=%d, EMSG=%s]", e.Code, msg)
	}
	return fmt.Sprintf("[Not Found]")
}

func NewError(code int32) TaskDevError {
	e := TaskDevError{}
	e.Code = code
	return e
}

func init() {
	err = make(map[int32]string)

	err[KESuccess] = "处理成功"

}
