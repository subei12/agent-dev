package telemetry

import "context"

type Provider struct{}

func NewProvider(_ string) *Provider {
	return &Provider{}
}

func (p *Provider) Shutdown(_ context.Context) error {
	return nil
}
