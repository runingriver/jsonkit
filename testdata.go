package jsonkit

import "github.com/jsonkit/conf"

var (
	unmarshalAllConfig = conf.AnnotateConfig{
		JsonStr: &conf.JsonStrAction{
			CheckAll:    true,
			CheckByPath: nil,
		},
	}
	unmarshalPathConfig = conf.AnnotateConfig{
		JsonStr: &conf.JsonStrAction{
			CheckAll:    false,
			CheckByPath: []string{"s.info", "s.content"},
		},
	}
	normalCaseConfig = conf.AnnotateConfig{
		JsonStr: &conf.JsonStrAction{
			CheckAll:    true,
			CheckByPath: nil,
		},
		NormalAnnotate: map[string]*conf.AnnotateFmt{
			"a.b.c": {
				Str: "测试注释",
			},
			"a.d.1.xx": {
				Str: "测试数组路径注释",
			},
		},
		EnumAnnotate: map[string]*conf.AnnotateFmt{
			"a.e.f": {UnitCvt: &conf.UnitCvt{
				Unit:   "元",
				Factor: 10,
				Op:     "/",
			}},
			"a.d.2.*": {EnumMap: map[string]string{"yy": "YY", "zz": "ZZ", "nn": "NN"}},
		},
		ArrayAnnotate: map[string]*conf.AnnotateFmt{
			"a.d.3.h": {EnumMap: map[string]string{"1": "一", "2": "二", "3": "三"}},
		},
	}
	normalConfigDataCase = `
	{
		"s": {
			"info": "{\"key1\":\"{\\\"nested_key1\\\":\\\"nested_value1\\\"}\"}",
			"content": "{\"7351241250965703962\":\"item_id\",\"2329\":\"app_id\",\"23.29\":\"high\"}",
			"item":"{\"item_id\":\"7351241250965703962\"}"
		},
		"a": {
			"b": {
				"c": "1"
			},
			"d": [
				{
					"xx": "yy"
				},
				{
					"xx": "zz"
				},
				{
					"xx": "yy",
					"mm": "nn"
				},
				{
					"h": [
						1,
						2,
						3
					]
				}
			],
			"e": {
				"f": "2"
			}
		}
	}	
	`
	// 特殊情况的case
	// 路径不合法;数字key不是index;非法的conf和data.
)

var (
	cutConfig = conf.CutterConfig{
		IncludeCut: []string{
			"name", "vendor.name", "vendor.names.1",
		},
		ExcludeCut: []string{
			"vendor.info.relation.Tom", "vendor.items.1", "vendor.prices.0",
		},
		DeepCut:    nil,
		OverrunCut: nil,
	}
	cutCaseOfStr = `
	{
		"name": "computers",
		"description": "List of computer products",
		"vendor": {
			"name": "Star Trek",
			"email": "info@example.com",
			"info": {
				"age": 25,
				"relation": {
					"Kite": "friend",
					"Tom": "son"
				}
			},
			"items": [
				{
					"id": 1,
					"name": "MacBook Pro 13 inch retina",
					"price": 1350
				},
				{
					"id": 2,
					"name": "MacBook Pro 15 inch retina",
					"price": 1700
				}
			],
			"prices": [
				2400,
				400.87,
				89.9,
				150.2
			],
			"names": [
				"John Doe",
				"Jane Doe",
				"Tom"
			]
		}
	}
	`
)
