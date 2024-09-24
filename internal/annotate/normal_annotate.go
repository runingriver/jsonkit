package annotate

import (
	"context"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/internal/pathconv"
	"github.com/jsonkit/jkerr"
	"github.com/jsonkit/jklog"

	"github.com/runingriver/mapinterface/mapitf"
)

var (
	NormalAnnotateIns = &NormalAnnotate{}
)

type NormalAnnotate struct {
}

func (e *NormalAnnotate) DoNormalAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	if len(jsonMap) == 0 || len(cfg) == 0 {
		return jsonMap, nil
	}
	for path, annotate := range cfg {
		_, err := e.DoOnePath(ctx, jsonMap, path, annotate)
		if err != nil {
			if !conf.SkipPartErr(ctx) {
				return jsonMap, err
			}
			jklog.CtxWarn(ctx, "DoNormalAnnotate exception,path:%s,err:%v", path, err)
		}
	}
	return jsonMap, nil
}

func (e *NormalAnnotate) DoOnePath(ctx context.Context, jsonMap map[string]interface{}, path string, fmt *conf.AnnotateFmt) (map[string]interface{}, error) {
	if err := pathconv.ValidNormalPath(path); err != nil || !conf.AnnotateFmtValid(fmt) {
		return nil, err
	}

	pathVal := mapitf.Fr(ctx, jsonMap)
	paths := pathconv.DotPathToSlice(path)
	endKey := paths[len(paths)-1]

	if len(paths) > 1 {
		for _, p := range paths[0 : len(paths)-1] {
			switch p.KeyType {
			case pathconv.NumKeyType:
				if _, ok := pathVal.New().Exist(p.Key); ok {
					pathVal = pathVal.Get(p.Key)
				} else {
					pathVal = pathVal.Index(p.Idx)
				}
			case pathconv.StrKeyType:
				pathVal = pathVal.Get(p.Key)
			case pathconv.AllKeyType:
				return jsonMap, jkerr.New(jkerr.PathIllegalErr, "un-support * path in NormalAnnotate Config")
			}
		}
	}

	valStr, err := pathVal.New().Get(endKey.Key).ToStr()
	if err != nil {
		return jsonMap, err
	}
	if v, ok := NotateVal(valStr, fmt); ok {
		_, err = pathVal.New().SetMap(endKey.Key, v)
		if err != nil {
			return nil, err
		}
	}

	return jsonMap, nil
}
