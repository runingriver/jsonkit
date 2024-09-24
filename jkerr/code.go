package jkerr

//go:generate stringer -type=JKErrorCode
type JKErrorCode int

const (
	UnknownErr       JKErrorCode = -1
	InitParseFailed  JKErrorCode = 10001
	InitParamTypeErr JKErrorCode = 20001
	ExceptObject     JKErrorCode = 20002

	PathIllegalErr JKErrorCode = 30001
	IterAllKeyErr  JKErrorCode = 30002
	ExcludePathErr JKErrorCode = 30003
)
