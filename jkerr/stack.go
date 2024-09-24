package jkerr

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// 堆栈的深度
const depth = 10

// stack 记录错误创建时的堆栈信息
type stacker []uintptr

// Location 返回格式:file_name:line func_name
func (s *stacker) Location() string {
	if len(*s) == 0 {
		return ""
	}
	frames := runtime.CallersFrames((*s)[:])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d %s", fileName(frame.File), frame.Line, funcName(frame.Function))
}

// Format 格式化输出,支持%s,%v,%+v,%#v
func (s *stacker) Format(st fmt.State, verb rune) {
	frames := runtime.CallersFrames((*s)[:])
	frame, b := frames.Next()
	switch verb {
	case 's':
		_, _ = fmt.Fprintf(st, "%s:%d %s", fileName(frame.File), frame.Line, funcName(frame.Function))
	case 'v':
		switch {
		case st.Flag('+'):
			for b {
				_, _ = io.WriteString(st, "\n")
				_, _ = fmt.Fprintf(st, "%s:%d", frame.File, frame.Line)
				frame, b = frames.Next()
			}
		case st.Flag('#'):
			for b {
				_, _ = io.WriteString(st, "\n")
				_, _ = fmt.Fprintf(st, "%s:%d %s", frame.File, frame.Line, funcName(frame.Function))
				frame, b = frames.Next()
			}
		default:
			for b {
				_, _ = io.WriteString(st, "\n")
				_, _ = io.WriteString(st, frame.Function)
				_, _ = io.WriteString(st, "\n\t")
				_, _ = fmt.Fprintf(st, "%s:%d", frame.File, frame.Line)
				frame, b = frames.Next()
			}
		}
	}
}

// funcName 将"xxx/xxx/common/berr.Hello"转换成:"Hello"
func funcName(frameName string) string {
	i := strings.LastIndex(frameName, "/")
	frameName = frameName[i+1:]
	i = strings.Index(frameName, ".")
	return frameName[i+1:]
}

func fileName(frameFile string) string {
	i := strings.LastIndex(frameFile, "/")
	return frameFile[i+1:]
}

func stack() *stacker {
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stacker = pcs[0:n]
	return &st
}
