// Tests for the jobs dequeuer.
package dequeuer

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/kevinburke/rickover/models"
)

type DummyProcessor struct {
	Count int64
}

func (dp *DummyProcessor) DoWork(_ *models.QueuedJob) error {
	atomic.AddInt64(&dp.Count, 1)
	return nil
}

func (dp *DummyProcessor) Sleep(_ uint32) time.Duration {
	return 0
}

type ChannelProcessor struct {
	Count int64
	Ch    chan struct{}
}

func (dp *ChannelProcessor) DoWork(qj *models.QueuedJob) error {
	select {
	case dp.Ch <- struct{}{}:
		atomic.AddInt64(&dp.Count, 1)
		return nil
	case <-time.After(100 * time.Millisecond):
		return errors.New("channel send timed out")
	}
}

func (dp *ChannelProcessor) Sleep(_ uint32) time.Duration {
	return 0
}
