package block

import (
	"time"

	"github.com/scaleway/scaleway-sdk-go/internal/async"
	"github.com/scaleway/scaleway-sdk-go/internal/errors"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// WaitForSnapshotRequest is used by WaitForSnapshot method.
type WaitForSnapshotRequest struct {
	SnapshotID    string
	Zone          scw.Zone
	Timeout       *time.Duration
	RetryInterval *time.Duration
}

// WaitForSnapshot wait for the snapshot to be in a "terminal state" before returning.
func (s *API) WaitForSnapshot(req *WaitForSnapshotRequest, opts ...scw.RequestOption) (*Snapshot, error) {
	timeout := defaultTimeout
	if req.Timeout != nil {
		timeout = *req.Timeout
	}
	retryInterval := defaultRetryInterval
	if req.RetryInterval != nil {
		retryInterval = *req.RetryInterval
	}

	terminalStatus := map[SnapshotStatus]struct{}{
		SnapshotStatusAvailable: {},
		SnapshotStatusInUse:     {},
		SnapshotStatusError:     {},
		SnapshotStatusLocked:    {},
		SnapshotStatusDeleted:   {},
	}

	snapshot, err := async.WaitSync(&async.WaitSyncConfig{
		Get: func() (interface{}, bool, error) {
			res, err := s.GetSnapshot(&GetSnapshotRequest{
				SnapshotID: req.SnapshotID,
				Zone:       req.Zone,
			}, opts...)
			if err != nil {
				return nil, false, err
			}
			_, isTerminal := terminalStatus[res.Status]

			return res, isTerminal, err
		},
		Timeout:          timeout,
		IntervalStrategy: async.LinearIntervalStrategy(retryInterval),
	})
	if err != nil {
		return nil, errors.Wrap(err, "waiting for snapshot failed")
	}
	return snapshot.(*Snapshot), nil
}
