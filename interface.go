package jsonkit

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/internal/annotate"
	"github.com/jsonkit/internal/cutter"
	"github.com/jsonkit/pkg/tool"
)

type JsonProcessOpt struct {
	Config string
}

type JsonProcessOption func(*JsonProcessOpt)

func WithConfig(config string) JsonProcessOption {
	return func(o *JsonProcessOpt) {
		o.Config = config
	}
}

// JsonProcess 传入tcc或byte conf配置进行json的cut和annotate
func JsonProcess(ctx context.Context, jsonMap map[string]interface{}, opts ...JsonProcessOption) (m map[string]interface{}, err error) {
	opt := &JsonProcessOpt{}
	for _, o := range opts {
		o(opt)
	}
	if opt.Config == "" {
		return jsonMap, errors.New("param illegal")
	}
	var actionCfg *conf.ActionConfig
	actionCfg, err = conf.LocalLoader(opt.Config)
	if err != nil {
		return nil, err
	}
	if actionCfg == nil || (actionCfg.AnnotateConfig == nil && actionCfg.CutterConfig == nil) {
		return jsonMap, nil
	}
	ctx = conf.CtxWithConf(ctx, actionCfg)
	annotateMap, err := JsonAnnotate(ctx, jsonMap, actionCfg.AnnotateConfig)
	if err != nil {
		return nil, err
	}
	cutMap, err := JsonMapCutter(ctx, annotateMap, actionCfg.CutterConfig)
	if err != nil {
		return nil, err
	}
	return cutMap, nil
}

// JsonAnnotate Json注释,给JsonVal上加注释
func JsonAnnotate(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.AnnotateConfig) (map[string]interface{}, error) {
	return annotate.GetJsonAnnotate().JsonAnnotate(ctx, jsonMap, cfg)
}

// JsonStrAnnotate Json注释,给JsonVal上加注释
func JsonStrAnnotate(ctx context.Context, jsonStr string, cfg *conf.AnnotateConfig) (string, error) {
	if len(jsonStr) == 0 || !json.Valid(tool.StrToByte(jsonStr)) || cfg == nil {
		return jsonStr, nil
	}
	jsonMap, err := tool.JsonLoadsMap(jsonStr)
	if err != nil {
		return "", err
	}
	processedMap, err := annotate.GetJsonAnnotate().JsonAnnotate(ctx, jsonMap, cfg)
	if err != nil {
		return "", err
	}
	return tool.JsonDumps(processedMap)
}

// JsonStrCutter 对json str进行剪裁
func JsonStrCutter(ctx context.Context, jsonStr string, cfg *conf.CutterConfig) (string, error) {
	if len(jsonStr) == 0 || !json.Valid(tool.StrToByte(jsonStr)) || cfg == nil {
		return jsonStr, nil
	}
	jsonMap, err := tool.JsonLoadsMap(jsonStr)
	if err != nil {
		return "", err
	}
	cutMap, err := cutter.GetJsonCutter().CutterMap(ctx, jsonMap, cfg)
	if err != nil {
		return "", err
	}
	return tool.JsonDumps(cutMap)
}

// JsonMapCutter 剪裁map
func JsonMapCutter(ctx context.Context, jsonMap map[string]interface{}, cfg *conf.CutterConfig) (map[string]interface{}, error) {
	return cutter.GetJsonCutter().CutterMap(ctx, jsonMap, cfg)
}

func GlobalConfig() *conf.GlobalConfig {
	return conf.Conf()
}
