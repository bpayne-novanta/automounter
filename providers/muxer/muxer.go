package muxer

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/pauldotknopf/automounter/providers"
)

type muxer struct {
	p []providers.MediaProvider
}

// Create a muxer from multiple providers
func Create(p []providers.MediaProvider) providers.MediaProvider {
	return &muxer{
		p,
	}
}

func (s *muxer) Initialize() error {
	for _, provider := range s.p {
		err := provider.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *muxer) Name() string {
	return "muxer"
}

func (s *muxer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var eg errgroup.Group
	for _, provider := range s.p {
		eg.Go(func() error {
			err := provider.Start(ctx)
			if err != nil {
				cancel()
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

func (s *muxer) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, provider := range s.p {
		result = append(result, provider.GetMedia()...)
	}
	return result
}

func (s *muxer) Mount(id string) (providers.MountSession, error) {
	for _, provider := range s.p {
		session, err := provider.Mount(id)
		if err == providers.ErrIDNotFound {
			continue
		}
		if err == nil {
			return session, nil
		}
		return nil, err
	}
	return nil, providers.ErrIDNotFound
}

func (s *muxer) Unmount(id string) error {
	for _, provider := range s.p {
		err := provider.Unmount(id)
		if err == providers.ErrIDNotFound {
			continue
		}
		return err
	}
	return providers.ErrIDNotFound
}
