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
	UnmarshalStrIns = &UnmarshalStr{}
)

type UnmarshalStr struct {
}

func (u *UnmarshalStr) UnmarshalStr(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.JsonStrAction) (map[string]interface{}, error) {
	if cfg.CheckAll {
		_, err := mapitf.Fr(ctx, jsonMap).SetAllAsMap()
		if err != nil {
			return jsonMap, err
		}
		return jsonMap, nil
	}

	if len(cfg.CheckByPath) == 0 {
		return jsonMap, nil
	}

	for _, path := range cfg.CheckByPath {
		pathVal := mapitf.Fr(ctx, jsonMap)
		paths := pathconv.DotPathToSlice(path)
		for _, p := range paths {
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
				return jsonMap, jkerr.New(jkerr.PathIllegalErr, "un-support * path in UnmarshalAction Config")
			}
		}
		_, err := pathVal.SetAllAsMap()
		if err != nil {
			if !conf.SkipPartErr(ctx) {
				return jsonMap, err
			}
			jklog.CtxWarn(ctx, "UnmarshalStr.CheckByPath SetAllAsMap err:%v,path:%s", err, path)
		}
	}

	return jsonMap, nil
}
