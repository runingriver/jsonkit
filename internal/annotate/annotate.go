package annotate

import (
	"context"
	"sync"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/jklog"
)

var (
	jsonAnnotate *JsonAnnotateImpl
	once         sync.Once
)

func GetJsonAnnotate() *JsonAnnotateImpl {
	if jsonAnnotate == nil {
		once.Do(func() {
			jsonAnnotate = &JsonAnnotateImpl{}
		})
	}

	return jsonAnnotate
}

type JsonAnnotateImpl struct {
}

func (j *JsonAnnotateImpl) JsonAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.AnnotateConfig) (map[string]interface{}, error) {
	if len(jsonMap) == 0 || cfg == nil {
		return jsonMap, nil
	}

	// 层层递归,把json str反序列化为map
	doUnmarshalMap, err := j.DoForUnmarshalStr(ctx, jsonMap, cfg.JsonStr)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "JsonAnnotate.DoForUnmarshalStr err:%v", err)
	}

	doNormalMap, err := j.DoForNormalAnnotate(ctx, doUnmarshalMap, cfg.NormalAnnotate)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "JsonAnnotate.DoForNormalAnnotate err:%v", err)
	}

	doEnumMap, err := j.DoForEnumAnnotate(ctx, doNormalMap, cfg.EnumAnnotate)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "JsonAnnotate.DoForEnumAnnotate err:%v", err)
	}

	doArrayMap, err := j.DoForArrayAnnotate(ctx, doEnumMap, cfg.ArrayAnnotate)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "JsonAnnotate.DoForArrayAnnotate err:%v", err)
	}

	return doArrayMap, nil
}

func (j *JsonAnnotateImpl) DoForUnmarshalStr(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.JsonStrAction) (map[string]interface{}, error) {
	return UnmarshalStrIns.UnmarshalStr(ctx, jsonMap, cfg)
}

func (j *JsonAnnotateImpl) DoForNormalAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	return NormalAnnotateIns.DoNormalAnnotate(ctx, jsonMap, cfg)
}

func (j *JsonAnnotateImpl) DoForEnumAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	return EnumAnnotateIns.DoEnumAnnotate(ctx, jsonMap, cfg)
}

func (j *JsonAnnotateImpl) DoForArrayAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	return ArrayAnnotateIns.DoArrayAnnotate(ctx, jsonMap, cfg)
}
