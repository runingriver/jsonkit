package tool

import (
	"errors"
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

// JsonLoadsMap
// Tip: map类型的json key必须是字符串;
func JsonLoadsMap(jsonStr string) (map[string]interface{}, error) {
	var m map[string]interface{}
	dc := decoder.NewDecoder(jsonStr)
	dc.UseNumber()

	if err := dc.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

// JsonLoadsList
func JsonLoadsList(jsonStr string) ([]interface{}, error) {
	var m []interface{}
	dc := decoder.NewDecoder(jsonStr)
	dc.UseNumber()

	if err := dc.Decode(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

// JsonLoadsObj 注:o必须是对象的指针
func JsonLoadsObj(jsonStr string, o interface{}) (interface{}, error) {
	if t := reflect.TypeOf(o); t.Kind() != reflect.Ptr {
		return nil, errors.New("param o must be ptr")
	}
	dc := decoder.NewDecoder(jsonStr)
	dc.UseNumber()

	if err := dc.Decode(o); err == nil {
		return o, nil
	} else {
		return o, err
	}
}

// JsonUnmarshalObj 注:o必须是对象的指针
func JsonUnmarshalObj(jsonStr string, o interface{}) (interface{}, error) {
	if t := reflect.TypeOf(o); t.Kind() != reflect.Ptr {
		return nil, errors.New("param o must be ptr")
	}
	err := sonic.Unmarshal(StrToByte(jsonStr), o)
	return o, err
}

func JsonDumps(m interface{}) (string, error) {
	if str, err := sonic.MarshalString(m); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(str), nil
	}
}

// StrChecker 判断v是不是一个字符串类型,首先判断是否为字符串,其次判断是否为json
func JsonChecker(v interface{}) (bool, string) {
	if isStr, s := IsStrType(v); isStr {
		if ok := sonic.Valid(StrToByte(s)); ok { // 底层调用的json.Valid
			return true, s
		}
		return false, s
	}

	return false, ""
}
