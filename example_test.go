package jsonkit

import (
	"context"
	"testing"

	"github.com/runingriver/mapinterface/mapitf"
	"github.com/runingriver/mapinterface/pkg"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalAllStr(t *testing.T) {
	ctx := context.Background()
	jsonMap, _ := pkg.JsonLoadsMap(normalConfigDataCase)
	annotate, err := JsonAnnotate(ctx, jsonMap, &unmarshalAllConfig)
	assert.Nil(t, err)
	assert.IsType(t, map[string]interface{}{}, annotate["s"].(map[string]interface{})["info"])
	assert.IsType(t, map[string]interface{}{}, annotate["s"].(map[string]interface{})["content"])
}

func TestUnmarshalPathStr(t *testing.T) {
	ctx := context.Background()
	jsonMap, _ := pkg.JsonLoadsMap(normalConfigDataCase)
	annotate, err := JsonAnnotate(ctx, jsonMap, &unmarshalPathConfig)
	assert.Nil(t, err)
	assert.IsType(t, map[string]interface{}{}, annotate["s"].(map[string]interface{})["info"])
	assert.IsType(t, map[string]interface{}{}, annotate["s"].(map[string]interface{})["content"])
	assert.IsType(t, "", annotate["s"].(map[string]interface{})["item"])
}

func TestNormalCase(t *testing.T) {
	ctx := context.Background()
	jsonMap, _ := pkg.JsonLoadsMap(normalConfigDataCase)
	annotate, err := JsonAnnotate(ctx, jsonMap, &normalCaseConfig)
	assert.Nil(t, err)

	val, _ := mapitf.From(annotate).GetAny("a", "b", "c").ToStr()
	assert.Equal(t, "1 (测试注释)", val)

	val, _ = mapitf.From(annotate).GetAny("a", "d").Index(1).Get("xx").ToStr()
	assert.Equal(t, "zz (测试数组路径注释)", val)

	val, _ = mapitf.From(annotate).GetAny("a", "e", "f").ToStr()
	assert.Equal(t, "2 (0.2 元)", val)

	val, _ = mapitf.From(annotate).GetAny("a", "d").Index(2).ToStr()
	// 这行可能失败,由于顺序问题
	assert.Equal(t, "{\"xx\":\"yy (YY)\",\"mm\":\"nn (NN)\"}", val)

	val, _ = mapitf.From(annotate).GetAny("a", "d").Index(3).Get("h").ToStr()
	assert.Equal(t, "[\"1 (一)\",\"2 (二)\",\"3 (三)\"]", val)
}

func TestCutterStr(t *testing.T) {
	ctx := context.Background()
	json, err := JsonStrCutter(ctx, cutCaseOfStr, &cutConfig)
	assert.Nil(t, err)
	assert.Equal(t, 235, len(json))
	t.Logf("%s,len:%d", json, len(json))
}
