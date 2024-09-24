package jkmetric

import "errors"

// MetricsClient 打点支持,不定义具体的实现,通过SetMetricsClient注入,注入的对象务必实现了该接口
type MetricsClient interface {
	EmitCounter(name string, value interface{}, tags map[string]string) error
	EmitTimer(name string, value interface{}, tags map[string]string) error
	EmitStore(name string, value interface{}, tags map[string]string) error
}

var (
	MetricsIns = &MetricsImpl{}
)

type MetricsImpl struct {
	metricsClient interface{}
}

func SetMetricsClient(client interface{}) MetricsClient {
	MetricsIns.metricsClient = client
	return MetricsIns
}

func (m *MetricsImpl) EmitCounter(name string, value interface{}, tags map[string]string) error {
	if m.metricsClient == nil {
		return errors.New("metrics client un-injected")
	}
	return m.metricsClient.(MetricsClient).EmitCounter(name, value, tags)
}

func (m *MetricsImpl) EmitTimer(name string, value interface{}, tags map[string]string) error {
	if m.metricsClient == nil {
		return errors.New("metrics client un-injected")
	}
	return m.metricsClient.(MetricsClient).EmitTimer(name, value, tags)
}

func (m *MetricsImpl) EmitStore(name string, value interface{}, tags map[string]string) error {
	if m.metricsClient == nil {
		return errors.New("metrics client un-injected")
	}
	return m.metricsClient.(MetricsClient).EmitStore(name, value, tags)
}
