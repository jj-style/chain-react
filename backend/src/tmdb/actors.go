package tmdb

import (
	"context"
	"fmt"
	"sync/atomic"

	go_tmdb "github.com/jj-style/go-tmdb"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (t *tmdb) getActorsBetween(ctx context.Context, c chan<- *go_tmdb.Person, r func() error, start, end int) error {
	g, ctx := errgroup.WithContext(ctx)

	ids := make(chan int)
	// Produce
	g.Go(func() error {
		defer close(ids)
		for i := start; i <= end; i++ {
			id := i // https://golang.org/doc/faq#closures_and_goroutines
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ids <- id:
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

			for id := range ids {
				if p, err := t.client.GetPersonInfo(id, map[string]string{"language": "en-GB"}); err != nil {
					t.log.Errorf("fetching person(%d): %v", id, err)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case c <- p:
					}
				}
			}
			return nil
		})
	}

	g.Go(r)

	return g.Wait()
}

func (t *tmdb) getLatestPerson() (*go_tmdb.PersonLatest, error) {
	latest, err := t.client.GetPersonLatest()
	if err != nil {
		return nil, err
	}
	t.log.WithFields(log.Fields{"name": latest.Name, "id": latest.ID}).Info("got latest person")
	return latest, nil
}

func (t *tmdb) GetAllActors(ctx context.Context, c chan<- *go_tmdb.Person, r func() error) error {
	latest, err := t.getLatestPerson()
	if err != nil {
		return err
	}

	latestId := latest.ID
	// latestId = 5
	return t.getActorsBetween(ctx, c, r, 1, latestId)
}

func (t *tmdb) GetActorsFrom(ctx context.Context, c chan<- *go_tmdb.Person, r func() error, id int) error {
	latest, err := t.getLatestPerson()
	if err != nil {
		return err
	}

	latestId := latest.ID
	//latestId = id + 5

	// sanity check
	if id >= latestId {
		return fmt.Errorf("Current latest id(%d) is greater than the latest from TMDB(%d)", id, latestId)
	}

	return t.getActorsBetween(ctx, c, r, id+1, latestId)
}
