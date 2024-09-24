package annotate

import (
	"context"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/internal/pathconv"
	"github.com/jsonkit/jkerr"
	"github.com/jsonkit/jklog"

	"github.com/runingriver/mapinterface/api"
	"github.com/runingriver/mapinterface/mapitf"
)

var (
	ArrayAnnotateIns = &ArrayAnnotate{}
)

type ArrayAnnotate struct {
}

func (e *ArrayAnnotate) DoArrayAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg map[string]*conf.AnnotateFmt) (map[string]interface{}, error) {
	if len(jsonMap) == 0 || len(cfg) == 0 {
		return jsonMap, nil
	}
	for path, fmt := range cfg {
		_, err := e.DoOnePath(ctx, jsonMap, path, fmt)
		if err != nil {
			jklog.CtxWarn(ctx, "DoEnumAnnotate exception,path:%s,err:%v", path, err)
			return jsonMap, err
		}
	}
	return jsonMap, nil
}

// DoOnePath x.y.z z对应的val一定是个list,z不能是index.
func (e *ArrayAnnotate) DoOnePath(ctx context.Context, jsonMap map[string]interface{}, path string, fmt *conf.AnnotateFmt) (m map[string]interface{}, err error) {
	if err := pathconv.ValidPath(path); err != nil || !conf.AnnotateFmtValid(fmt) {
		return nil, err
	}

	pathVal := mapitf.Fr(ctx, jsonMap)
	paths := pathconv.DotPathToSlice(path)
	var prePathVal api.MapInterface
	for i, p := range paths {
		if i == len(paths)-1 {
			prePathVal = pathVal.New()
		}
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

	strList, err := pathVal.New().ToListStrF()
	if err != nil {
		return nil, err
	}
	resultList := make([]string, 0, len(strList))
	for _, s := range strList {
		v, _ := NotateVal(s, fmt)
		resultList = append(resultList, v)
	}

	endKey := paths[len(paths)-1]
	_, err = prePathVal.SetMap(endKey.Key, resultList)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
