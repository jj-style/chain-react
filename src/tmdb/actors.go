package tmdb

import (
	"context"
	"fmt"
	"sync/atomic"

	go_tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (t *tmdb) GetAllActors(ctx context.Context) (map[int]*go_tmdb.Person, error) {
	latest, err := t.client.GetPersonLatest()
	if err != nil {
		return map[int]*go_tmdb.Person{}, nil
	}
	t.log.WithFields(log.Fields{"name": latest.Name, "id": latest.ID}).Info("got latest person")

	g, ctx := errgroup.WithContext(ctx)

	ids := make(chan int)
	// Produce
	g.Go(func() error {
		defer close(ids)
		latestId := 10
		//latestId := latest.ID
		for i := 1; i <= latestId; i++ {
			id := i // https://golang.org/doc/faq#closures_and_goroutines
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ids <- id:
			}
		}
		return nil
	})

	people := make(chan *go_tmdb.Person)
	// Map
	workers := int32(10)
	for i := 0; i < int(workers); i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(people)
				}
			}()

			for id := range ids {
				if p, err := t.client.GetPersonInfo(id, map[string]string{"language": "en"}); err != nil {
					return fmt.Errorf("GetPersonInfo %d: %s", id, err)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case people <- p:
					}
				}
			}
			return nil
		})
	}

	// Reduce
	ret := map[int]*go_tmdb.Person{}
	g.Go(func() error {
		for p := range people {
			ret[p.ID] = p
		}
		return nil
	})
	return ret, g.Wait()
}
