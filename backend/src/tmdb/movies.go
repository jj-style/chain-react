package tmdb

import (
	"context"
	"sync/atomic"

	go_tmdb "github.com/jj-style/go-tmdb"
	"golang.org/x/sync/errgroup"
)

func (t *tmdb) GetActorMovieCredits(ctx context.Context, id int) (*go_tmdb.PersonMovieCredits, error) {
	return t.client.GetPersonMovieCredits(id, map[string]string{"language": "en-GB"})
}

func (t *tmdb) GetAllActorMovieCredits(ctx context.Context, c chan<- *go_tmdb.PersonMovieCredits, r func() error, ids ...int) error {
	g, ctx := errgroup.WithContext(ctx)

	actorIds := make(chan int)
	// Produce
	g.Go(func() error {
		defer close(actorIds)
		for _, i := range ids {
			id := i // https://golang.org/doc/faq#closures_and_goroutines
			select {
			case <-ctx.Done():
				return ctx.Err()
			case actorIds <- id:
			}
		}
		return nil
	})

	// Map
	workers := int32(10)
	for i := 0; i < int(workers); i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(c)
				}
			}()

			for id := range actorIds {
				if credits, err := t.GetActorMovieCredits(ctx, id); err != nil {
					t.log.Errorf("fetching credits for person(%d): %v", id, err)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case c <- credits:
					}
				}
			}
			return nil
		})
	}

	g.Go(r)

	return g.Wait()
}
