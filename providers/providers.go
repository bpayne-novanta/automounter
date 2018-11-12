package providers

import (
	"context"
)

var providers = struct {
	p []MediaProvider
}{}

// MediaProvider The type that will detect and mount media
type MediaProvider interface {
	Name() string
	Start(context.Context) error
	GetMedia() []Media
}

// Media A media type that can be mounted/used.
type Media interface {
	ID() string
	DisplayName() string
}

// AddProvider Add a provider
func AddProvider(provider MediaProvider) {
	providers.p = append(providers.p, provider)
}

// GetProviders Get the current providers
func GetProviders() []MediaProvider {
	return providers.p
}
