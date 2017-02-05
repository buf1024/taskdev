package util

import (
	"fmt"
)

const (
	E_SUCCESS uint32 = 99999 - iota
)

type TaskDevError struct {
	Code uint32
}

var err map[uint32]string

func (e TaskDevError) Error() string {
	if msg, ok := err[e.Code]; ok {
		return fmt.Sprintf("[ERR=%d, EMSG=%s]", e.Code, msg)
	}
	return fmt.Sprintf("[Not Found]")
}

func NewError(code uint32) TaskDevError {
	e := TaskDevError{}
	e.Code = code
	return e
}

func init() {
	err = make(map[uint32]string)

	err[E_SUCCESS] = "处理成功"

}
