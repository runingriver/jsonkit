package jkerr

import (
	"errors"
	"fmt"
	"io"
)

type JKErr interface {
	error
	Is(error) bool
	Wrap(error) JKErr
	ErrCode() JKErrorCode
	String() string
	IsErrEqual(err error) bool
}

type JKError struct {
	Code JKErrorCode
	Err  error
	Msg  string

	Stack *stacker
}

func New(code JKErrorCode, format string, args ...interface{}) JKErr {
	return &JKError{
		Code:  code,
		Msg:   fmt.Sprintf(format, args...),
		Stack: stack(),
	}
}

func (j *JKError) String() string {
	if j.Err == nil && j.Msg == "" {
		return fmt.Sprintf("RecordErr{Code:%v,Loc:%s}", j.Code, j.Stack.Location())
	}
	return fmt.Sprintf("RecordErr{Code:%d,Loc:%s,Err:%v,Msg:%s}", j.Code, j.Stack.Location(), j.Err, j.Msg)
}

func (j *JKError) Error() string {
	return j.String()
}

func (j *JKError) Wrap(err error) JKErr {
	j.Err = err
	return j
}

func (j *JKError) Is(err error) bool {
	return IsItfErr(err)
}

func (j *JKError) ErrCode() JKErrorCode {
	return j.Code
}

func (j *JKError) IsErrEqual(err2 error) bool {
	return IsErrEqual(j, err2)
}

func (j *JKError) Format(st fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = io.WriteString(st, j.String())
	case 'v':
		switch {
		case st.Flag('+'):
			_, _ = io.WriteString(st, j.String())
			_, _ = fmt.Fprintf(st, "%+v", j.Stack)
		case st.Flag('#'):
			_, _ = io.WriteString(st, j.String())
			_, _ = fmt.Fprintf(st, "%#v", j.Stack)
		default:
			_, _ = io.WriteString(st, j.String())
			_, _ = fmt.Fprintf(st, "%v", j.Stack)
		}
	}
}

func IsErrEqual(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return false
	}

	if errors.Is(err1, err2) {
		return true
	}

	err1Code := GetErrCode(err1)
	err2Code := GetErrCode(err2)
	if err1Code == err2Code && err1Code != -1 {
		return true
	}

	return false
}

func IsItfErr(err error) bool {
	if _, ok := err.(*JKError); ok {
		return true
	}
	return false
}

func GetErrCode(err error) JKErrorCode {
	if jkErr, ok := err.(*JKError); ok {
		return jkErr.Code
	}
	return UnknownErr
}
