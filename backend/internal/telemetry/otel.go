package telemetry

import "context"

type Provider struct{}

// NewProvider 创建并返回对应的组件。
func NewProvider(_ string) *Provider {
	return &Provider{}
}

// Shutdown 实现当前函数行为。
func (p *Provider) Shutdown(_ context.Context) error {
	return nil
}
