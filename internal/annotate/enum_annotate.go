package annotate

import (
	"context"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/internal/pathconv"
	"github.com/jsonkit/jkerr"
	"github.com/jsonkit/jklog"
	"github.com/jsonkit/pkg/tool"

	"github.com/runingriver/mapinterface/mapitf"
)

var (
	EnumAnnotateIns = &EnumAnnotate{}
)

type EnumAnnotate struct {
}

func (e *EnumAnnotate) DoEnumAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	if len(jsonMap) == 0 || len(cfg) == 0 {
		return jsonMap, nil
	}
	for path, annotate := range cfg {
		_, err := e.DoOnePath(ctx, jsonMap, path, annotate)
		if err != nil {
			if !conf.SkipPartErr(ctx) {
				return jsonMap, err
			}
			jklog.CtxWarn(ctx, "DoEnumAnnotate exception,path:%s,err:%v", path, err)
		}
	}
	return jsonMap, nil
}

func (e *EnumAnnotate) DoOnePath(ctx context.Context, jsonMap map[string]interface{}, path string, fmt *conf.AnnotateFmt) (m map[string]interface{}, err error) {
	if err := pathconv.ValidPath(path); err != nil || !conf.AnnotateFmtValid(fmt) {
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
				return jsonMap, jkerr.New(jkerr.PathIllegalErr, "un-support * path in EnumAnnotate Config")
			}
		}
	}

	// *结尾的场景: x.y.*, 最后一层所有key去匹配annotateMap中的key
	if endKey.KeyType == pathconv.AllKeyType {
		val, err := pathVal.New().ForEach(func(i int, k, v interface{}) (key, val interface{}) {
			mapKeyStr := tool.ToStr(k)
			if v, ok := NotateVal(v, fmt); ok {
				_, err := pathVal.New().SetMap(mapKeyStr, v)
				if err != nil {
					return nil, nil
				}
			}
			return k, v
		}).ToMapItf()
		if err != nil {
			return nil, err
		}

		cvtMap, _ := pathVal.New().ToMap()
		if len(val) != len(cvtMap) {
			return nil, jkerr.New(jkerr.IterAllKeyErr, "iter path of * exception,err:%v", err)
		}

		return jsonMap, err
	}

	// 正常场景
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
