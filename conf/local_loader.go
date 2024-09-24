package conf

import (
	"errors"

	"github.com/jsonkit/pkg/tool"

	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

func LocalLoader(configStr string) (*ActionConfig, error) {
	if len(configStr) == 0 {
		return nil, errors.New("config str is empty")
	}
	var actionConfig ActionConfig
	if err := sonic.UnmarshalString(configStr, &actionConfig); err == nil {
		return &actionConfig, nil
	}
	err := yaml.Unmarshal(tool.StrToByte(configStr), &actionConfig)
	return &actionConfig, err
}
