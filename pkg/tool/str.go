package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func IsStrType(v interface{}) (bool, string) {
	v = Interpret(v)
	if v == nil {
		return false, ""
	}

	switch vv := v.(type) {
	case string:
		return true, vv
	case []byte:
		return true, ByteToStr(vv)
	default:
		return false, ""
	}
}

func ToStr(v interface{}) string {
	result := ""
	if v == nil {
		return result
	}
	switch vv := v.(type) {
	case json.Number:
		result = vv.String()
	case string:
		result = vv
	case int:
		return strconv.FormatInt(int64(vv), 10)
	case int8:
		return strconv.FormatInt(int64(vv), 10)
	case int16:
		return strconv.FormatInt(int64(vv), 10)
	case int32:
		return strconv.FormatInt(int64(vv), 10)
	case int64:
		return strconv.FormatInt(vv, 10)
	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case uint64:
		return strconv.FormatUint(vv, 10)
	case float32:
		return strconv.FormatFloat(float64(vv), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(vv, 'G', -1, 64)
	case bool:
		result = strconv.FormatBool(vv)
	case []byte:
		result = ByteToStr(vv)
	case nil:
		return ``
	case error:
		return vv.Error()
	default:
		if f, ok := v.(fmt.Stringer); ok {
			return f.String()
		}

		if callRst, ok := CallMethod(v, "String"); ok {
			if result, ok = callRst.(string); ok {
				return result
			}
		}

		result, _ = JsonDumps(vv)
	}
	return result
}

func CallMethod(v interface{}, methodName string) (interface{}, bool) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(v)

	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		// 如果v是值类型,需要组装出一个指针类型的变量赋值给ptr
		ptr = reflect.New(reflect.TypeOf(v))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// 检查methodName是否存在于值类型的方法中
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// 检查methodName是否存在于指针类型的方法中
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		return finalMethod.Call([]reflect.Value{})[0].Interface(), true
	}

	// 无对应的方法实现
	return "", false
}

func Interpret(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// StrToByte 高效转换,避免内存拷贝
func StrToByte(s string) (b []byte) {
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return
}

// ByteToStr 高效转换,避免内存拷贝,
func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
