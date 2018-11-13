package leaser

import (
	"fmt"
	"sync"

	"github.com/pauldotknopf/automounter/helpers"
	"github.com/pauldotknopf/automounter/providers"
)

// Lease represents a leased media item
type Lease interface {
	ID() string
	MediaID() string
	MountPath() string
}

// Leaser The type that manages leases for media items
type Leaser interface {
	MediaProvider() providers.MediaProvider
	Leases() []Lease
	Lease(mediaID string) (Lease, error)
	Release(leaseID string) error
}

type leaser struct {
	mediaProvider providers.MediaProvider
	media         []*mediaLease
	lock          sync.Mutex
}

// Create a leaser object
func Create(mediaProvider providers.MediaProvider) Leaser {
	l := &leaser{}
	l.mediaProvider = mediaProvider
	l.media = make([]*mediaLease, 0)
	return l
}

func (s *leaser) MediaProvider() providers.MediaProvider {
	return s.mediaProvider
}

func (s *leaser) Leases() []Lease {
	s.lock.Lock()
	defer s.lock.Unlock()

	result := make([]Lease, 0)
	for _, media := range s.media {
		for _, lease := range media.leases {
			result = append(result, lease)
		}
	}
	return result
}

func (s *leaser) Lease(mediaID string) (Lease, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Look for an existing mount for this media item.
	for _, media := range s.media {
		if media.mediaID == mediaID {
			// This item currently is mounted, just add a lease.
			lease := &mediaLeaseItem{}
			lease.media = media
			lease.leaseID = helpers.RandString(10)
			media.leases = append(media.leases, lease)
			return lease, nil
		}
	}

	// This is lease for a new media item.
	mountSession, err := s.mediaProvider.Mount(mediaID)
	if err != nil {
		return nil, err
	}

	media := &mediaLease{}
	media.mediaID = mediaID
	media.MountSession = mountSession
	s.media = append(s.media, media)

	// Add one lease to the media item.
	lease := &mediaLeaseItem{}
	lease.media = media
	lease.leaseID = helpers.RandString(10)
	media.leases = append(media.leases, lease)

	return lease, nil
}

func (s *leaser) Release(leaseID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Look for an existing mount for this media item.
	for _, media := range s.media {
		for leaseIndex, lease := range media.leases {
			if lease.ID() == leaseID {
				media.leases = append(media.leases[:leaseIndex], media.leases[leaseIndex+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("no lease with the given id")
}
