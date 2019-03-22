package http

import (
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
	pR, pW := io.Pipe()

	// Consume reader.
	group.Go(func() error {
		s.Chainer.Read(ctx, pR)
		return nil
	})
	// Fan out.
	for i := 0; i < n; i++ {
		sem.Acquire(ctx, 1)
		group.Go(func() error {
			defer sem.Release(1)
			j, err := s.Joker.GetJoke(ctx)
			if err != nil {
				s.Logger.Errorw("GetJoke failed", "err", err)
			}
			// Fan in.
			pW.Write([]byte(fmt.Sprintf("%s\n", j)))
			return nil
		})
	}
	pW.Close()

	if err := group.Wait(); err != nil {
		return err
	}
	s.Logger.Infow("warmup completet", "count", n)
	return nil
}
