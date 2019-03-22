package http

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// WarmUp the Chainer by feeding it jokes from Joker.
func (s *Server) WarmUp(ctx context.Context, n int) error {
	sem := semaphore.NewWeighted(s.Config.MaxConcurrency)
	group, ctx := errgroup.WithContext(ctx)
	jokes := new(bytes.Buffer)
	pipeR, pipeW := io.Pipe()
	defer pipeR.Close()

	// Consume reader/channel
	group.Go(func() error {
		s.Chainer.Read(ctx, pipeR)
		return nil
	})

	for i := 0; i < n; i++ {
		defer pipeW.Close()
		sem.Acquire(ctx, 1)
		group.Go(func() error {
			defer sem.Release(1)
			j, err := s.Joker.GetJoke(ctx)
			if err != nil {
				s.Logger.Errorw("GetJoke failed", "err", err)
			}
			jokes.WriteString(fmt.Sprintf("%s\n", j))
			return nil
		})
	}
	return nil
}
