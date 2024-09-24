package cutter

import (
	"context"
	"sync"

	"github.com/Jeffail/gabs/v2"
	"github.com/jsonkit/conf"
	"github.com/jsonkit/internal/pathconv"
	"github.com/jsonkit/jkerr"
	"github.com/jsonkit/jklog"
)

var (
	jsonCutter *JsonCutterImpl
	once       sync.Once
)

func GetJsonCutter() *JsonCutterImpl {
	if jsonCutter == nil {
		once.Do(func() {
			jsonCutter = &JsonCutterImpl{}
		})
	}

	return jsonCutter
}

type JsonCutterImpl struct {
}

func (j *JsonCutterImpl) CutterMap(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.CutterConfig) (map[string]interface{}, error) {
	if len(jsonMap) == 0 || cfg == nil {
		return jsonMap, nil
	}

	includeMap, err := j.DoForIncludeCut(ctx, jsonMap, cfg.IncludeCut)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "CutterMap.DoForIncludeCut err:%v", err)
	}

	excludeMap, err := j.DoForExcludeCut(ctx, includeMap, cfg.ExcludeCut)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "CutterMap.DoForExcludeCut err:%v", err)
	}

	deepMap, err := j.DoForDeepCut(ctx, excludeMap, cfg.DeepCut)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "CutterMap.DoForDeepCut err:%v", err)
	}

	overrunMap, err := j.DoForOverrunCut(ctx, deepMap, cfg.OverrunCut)
	if err != nil {
		if !conf.SkipPartErr(ctx) {
			return jsonMap, err
		}
		jklog.CtxWarn(ctx, "CutterMap.DoForOverrunCut err:%v", err)
	}
	return overrunMap, nil
}

func (j *JsonCutterImpl) DoForIncludeCut(ctx context.Context, jsonMap map[string]interface{}, paths []string) (map[string]interface{}, error) {
	if !pathconv.ValidCutPath(paths) || !pathconv.ValidIncludeCutPath(paths) {
		return jsonMap, nil
	}
	for _, path := range paths {
		err := gabs.Wrap(jsonMap).DeleteP(path)
		if err != nil {
			if !conf.SkipPartErr(ctx) {
				return jsonMap, err
			}
			jklog.CtxWarn(ctx, "DoForIncludeCut del path err:%v,path:%v", err, path)
		}
	}
	return jsonMap, nil
}

func (j *JsonCutterImpl) DoForExcludeCut(ctx context.Context, jsonMap map[string]interface{}, paths []string) (map[string]interface{}, error) {
	if !pathconv.ValidCutPath(paths) || !pathconv.ValidExcludeCutPath(paths) {
		return jsonMap, nil
	}
	for _, path := range paths {
		preEnd, end := pathconv.PathToList(path)
		jsonObj := gabs.Wrap(jsonMap).S(preEnd...)

		switch vv := jsonObj.Data().(type) {
		case map[string]interface{}:
			for k := range vv {
				if end == k {
					continue
				}
				err := jsonObj.Delete(k)
				if err != nil {
					if !conf.SkipPartErr(ctx) {
						return jsonMap, err
					}
					jklog.CtxWarn(ctx, "DoForExcludeCut del map,err:%v,path:%v", err, path)
				}
			}
		case []interface{}:
			data := jsonObj.S(end).Data()
			_, err := gabs.Wrap(jsonMap).Set([]interface{}{data}, preEnd...)
			if err != nil {
				if !conf.SkipPartErr(ctx) {
					return jsonMap, err
				}
				jklog.CtxWarn(ctx, "DoForExcludeCut del list,err:%v,path:%v", err, path)
			}
		default:
			return nil, jkerr.New(jkerr.ExcludePathErr, "DoForExcludeCut path un-match data type:%s", path)
		}
	}
	return jsonMap, nil
}

func (j *JsonCutterImpl) DoForDeepCut(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.DeepCut) (map[string]interface{}, error) {
	return jsonMap, nil
}

func (j *JsonCutterImpl) DoForOverrunCut(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.OverrunCut) (map[string]interface{}, error) {
	return jsonMap, nil
}
