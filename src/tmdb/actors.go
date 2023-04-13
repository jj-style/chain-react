package tmdb

import (
	"fmt"

	go_tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (t *tmdb) GetAllActors() ([]*go_tmdb.Person, error) {
	latest, err := t.client.GetPersonLatest()
	if err != nil {
		return []*go_tmdb.Person{}, nil
	}
	t.log.WithFields(log.Fields{"name": latest.Name, "id": latest.ID}).Info("got latest person")

	g := new(errgroup.Group)

	latestId := 10
	//latestId := latest.ID
	people := make([]*go_tmdb.Person, 0, latestId)
	for i := 1; i <= latestId; i++ {
		id := i // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			p, err := t.client.GetPersonInfo(id, map[string]string{"language": "en"})
			if err != nil {
				return fmt.Errorf("fetching person(%d): %v", id, err)
			}
			people = append(people, p)
			return nil
		})
	}

	// Wait for all fetches to complete.
	if err := g.Wait(); err == nil {
		t.log.Info("successfully fetched all people")
	} else {
		t.log.Errorf("error fetching person: %v", err)
		return people, err
	}

	return people, nil
}
