package conf

import "context"

var (
	jsonKitConfig   = &GlobalConfig{}
	actionConfigKey = "jsonkit_config"
	skipPartErr     = true
)

func Conf() *GlobalConfig {
	return jsonKitConfig
}

type GlobalConfig struct {
	SkipPartErr *bool `json:"skip_part_err,omitempty"`
}

func (g *GlobalConfig) WithSkipPartErr() *GlobalConfig {
	g.SkipPartErr = &skipPartErr
	return g
}

type ActionConfig struct {
	SkipPartErr *bool `json:"skip_part_err,omitempty"`

	CutterConfig   *CutterConfig   `json:"cutter_config,omitempty" json:"cutter_config,omitempty"`
	AnnotateConfig *AnnotateConfig `json:"annotate_config,omitempty" json:"annotate_config,omitempty"`
}

func (a *ActionConfig) WithSkipPartErr() *ActionConfig {
	a.SkipPartErr = &skipPartErr
	return a
}

type CutterConfig struct {
	IncludeCut []string    `json:"include_cut,omitempty"`
	ExcludeCut []string    `json:"exclude_cut,omitempty"`
	DeepCut    *DeepCut    `json:"deep_cut,omitempty"`
	OverrunCut *OverrunCut `json:"overrun_cut,omitempty"`
}

type DeepCut struct {
	Paths []string `json:"paths,omitempty"`
	Deep  int      `json:"deep,omitempty"`
}

type OverrunCut struct {
	Paths      []string `json:"paths,omitempty"`
	LimitBytes int      `json:"limit_bytes,omitempty"`
}

// AnnotateConfig 支持的注释类型.
// 1. UnmarshalStr 对内嵌json str的处理行为
// 2. NormalAnnotate 直接在val上增加注释; 路径语法:a.b.c  a.b.*.c
// 3. EnumAnnotate 根据val值进行匹配; 路径语法:a.b.c  a.b.*.c
// 4. ArrayAnnotate val为数组,为每个数组值匹配; 路径语法:a.b.*.c.*
// 路径语法:a.b.c  a.b.*.c  a.b.*
type AnnotateConfig struct {
	JsonStr        *JsonStrAction          `json:"json_str,omitempty"`
	NormalAnnotate map[string]*AnnotateFmt `json:"normal_annotate,omitempty"`
	EnumAnnotate   map[string]*AnnotateFmt `json:"enum_annotate,omitempty"`
	ArrayAnnotate  map[string]*AnnotateFmt `json:"array_annotate,omitempty"`
}

type JsonStrAction struct {
	CheckAll    bool     `json:"check_all"`     // 所有字符串都检查并反序列化
	CheckByPath []string `json:"check_by_path"` // 检测指定路径下的字符串
}

// AnnotateFmt 按照Str->EnumMap->UnitOp的优先级逐步解析
type AnnotateFmt struct {
	Str string `json:"str,omitempty"` // 字符串注释

	UnitCvt *UnitCvt `json:"unit_cvt,omitempty"` // 单位换算

	EnumMap map[string]string `json:"enum_map,omitempty"` // 枚举map
}

type UnitCvt struct {
	Unit   string  `json:"unit,omitempty"`   // 单位转换注释, 单位,如:元
	Factor float64 `json:"factor,omitempty"` // 单位转换注释,算法因子,如:分转成元 = val/Op
	Op     string  `json:"op,omitempty"`     // 单位转换的算法,如:+,-,*,/
}

func AnnotateFmtValid(fmt *AnnotateFmt) bool {
	if fmt == nil {
		return false
	}
	if fmt.Str == "" && len(fmt.EnumMap) == 0 && fmt.UnitCvt == nil {
		return false
	}
	if fmt.UnitCvt != nil && (fmt.UnitCvt.Unit == "" || fmt.UnitCvt.Op == "" || fmt.UnitCvt.Factor == 0) {
		return false
	}
	return true
}

func CtxWithConf(ctx context.Context, ac *ActionConfig) context.Context {
	return context.WithValue(ctx, actionConfigKey, ac)
}

func SkipPartErr(ctx context.Context) bool {
	acObj := ctx.Value(actionConfigKey)
	if ac, ok := acObj.(*ActionConfig); ok {
		if ac.SkipPartErr != nil {
			return *ac.SkipPartErr
		}
	}
	if jsonKitConfig.SkipPartErr != nil {
		return *jsonKitConfig.SkipPartErr
	}
	return false
}
