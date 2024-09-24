package annotate

import (
	"fmt"

	"github.com/jsonkit/conf"
	"github.com/jsonkit/pkg/tool"

	"github.com/runingriver/mapinterface/pkg"
	"github.com/shopspring/decimal"
)

// NotateVal 给val加注释
func NotateVal(val interface{}, annotate *conf.AnnotateFmt) (string, bool) {
	if annotate.Str != "" {
		return fmt.Sprintf("%s (%s)", tool.ToStr(val), annotate.Str), true
	}
	if annotate.UnitCvt != nil {
		return UnitToStr(val, annotate.UnitCvt), true
	}

	if len(annotate.EnumMap) != 0 {
		valStr := tool.ToStr(val)
		if v, ok := annotate.EnumMap[valStr]; ok {
			return fmt.Sprintf("%s (%s)", valStr, v), true
		}
	}

	return tool.ToStr(val), false
}

func UnitToStr(val interface{}, unitCvt *conf.UnitCvt) string {
	v, err := pkg.ToInt64(val)
	if err != nil {
		return pkg.ToStr(val)
	}

	var num string
	switch unitCvt.Op {
	case "*":
		num = decimal.NewFromInt(v).Mul(decimal.NewFromFloat(unitCvt.Factor)).String()
	case "/":
		num = decimal.NewFromInt(v).DivRound(decimal.NewFromFloat(unitCvt.Factor), 2).String()
	case "+":
		num = decimal.NewFromInt(v).Add(decimal.NewFromFloat(unitCvt.Factor)).String()
	case "-":
		num = decimal.NewFromInt(v).Sub(decimal.NewFromFloat(unitCvt.Factor)).String()
	}

	if num != "" {
		return fmt.Sprintf("%s (%s %s)", tool.ToStr(val), num, unitCvt.Unit)
	}
	return pkg.ToStr(val)
}
