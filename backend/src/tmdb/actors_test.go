package tmdb_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jj-style/chain-react/src/tmdb"
	tmdbMocks "github.com/jj-style/chain-react/src/tmdb/mocks"
	go_tmdb "github.com/jj-style/go-tmdb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllActors(t *testing.T) {
	type args struct {
		called          int
		latestPerson    *go_tmdb.PersonLatest
		latestPersonErr error
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*tmdbMocks.MockTMDbClient, args)
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "All 10 Actors Happy",
			args: args{
				called:          10,
				latestPerson:    &go_tmdb.PersonLatest{ID: 10, Name: "Latest"},
				latestPersonErr: nil,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
				mtc.EXPECT().GetPersonInfo(mock.AnythingOfType("int"), mock.Anything).Return(&go_tmdb.Person{ID: 1, Name: "Person"}, nil).Times(args.called)
			},
			errFunc: assert.NoError,
		},
		{
			name: "All Actors - Unhappy Person Info Error No Person Received In Callback",
			args: args{
				called:          9,
				latestPerson:    &go_tmdb.PersonLatest{ID: 10, Name: "Latest"},
				latestPersonErr: nil,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
				for i := 1; i < 10; i++ {
					mtc.EXPECT().GetPersonInfo(i, mock.Anything).Return(&go_tmdb.Person{ID: i, Name: fmt.Sprintf("Person %d", i)}, nil).Once()
				}
				mtc.EXPECT().GetPersonInfo(10, mock.Anything).Return(nil, errors.New("boom")).Once()
			},
			errFunc: assert.NoError,
		},
		{
			name: "All Actors - Unhappy Latest Person Info Error",
			args: args{
				called:          0,
				latestPerson:    nil,
				latestPersonErr: errors.New("boom"),
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
			},
			errFunc: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			m := tmdbMocks.NewMockTMDbClient(t)
			if tt.setup != nil {
				tt.setup(m, tt.args)
			}
			tmdb := tmdb.NewClient(m, logrus.StandardLogger())

			c := make(chan *go_tmdb.Person)
			called := 0
			cb := func() error {
				for range c {
					called += 1
				}
				return nil
			}

			// Act
			err := tmdb.GetAllActors(context.TODO(), c, cb)

			// Assert
			tt.errFunc(t, err)
			assert.Equal(t, tt.args.called, called)
		})
	}

}

func TestGetActorsFrom(t *testing.T) {
	type args struct {
		called          int
		from            int
		latestPerson    *go_tmdb.PersonLatest
		latestPersonErr error
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*tmdbMocks.MockTMDbClient, args)
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "Get Actors From - Happy",
			args: args{
				called:          3,
				from:            7,
				latestPerson:    &go_tmdb.PersonLatest{ID: 10, Name: "Latest"},
				latestPersonErr: nil,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
				mtc.EXPECT().GetPersonInfo(mock.AnythingOfType("int"), mock.Anything).Return(&go_tmdb.Person{ID: 1, Name: "Person"}, nil).Times(args.called)
			},
			errFunc: assert.NoError,
		},
		{
			name: "Get Actors From - Unhappy from > latest",
			args: args{
				from:            5,
				latestPerson:    &go_tmdb.PersonLatest{ID: 3, Name: "Latest"},
				latestPersonErr: nil,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
			},
			errFunc: assert.Error,
		},
		{
			name: "Get Actors From - Unhappy Get Latest Actor Error",
			args: args{
				from:            5,
				latestPerson:    nil,
				latestPersonErr: errors.New("boom"),
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonLatest().Return(args.latestPerson, args.latestPersonErr)
			},
			errFunc: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			m := tmdbMocks.NewMockTMDbClient(t)
			if tt.setup != nil {
				tt.setup(m, tt.args)
			}
			tmdb := tmdb.NewClient(m, logrus.StandardLogger())

			c := make(chan *go_tmdb.Person)
			called := 0
			cb := func() error {
				for range c {
					called += 1
				}
				return nil
			}

			// Act
			err := tmdb.GetActorsFrom(context.TODO(), c, cb, tt.args.from)

			// Assert
			tt.errFunc(t, err)
			assert.Equal(t, tt.args.called, called)
		})
	}

}

func TestGetActorsByName(t *testing.T) {
	type args struct {
		called int
		names  []string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*tmdbMocks.MockTMDbClient, args)
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "Get Actors By Name - Happy",
			args: args{
				called: 4,
				names:  []string{"john", "paul", "george", "ringo"},
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				for _, name := range args.names {
					mtc.EXPECT().SearchPerson(name, mock.Anything).Return(&go_tmdb.PersonSearchResults{TotalResults: 1, Results: getSearchResults(name)}, nil)
				}
			},
			errFunc: assert.NoError,
		},
		{
			name: "Get Actors By Name - Unhappy Error Searching Person",
			args: args{
				called: 3,
				names:  []string{"john", "paul", "george", "ringo"},
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				for _, name := range args.names[:3] {
					mtc.EXPECT().SearchPerson(name, mock.Anything).Return(&go_tmdb.PersonSearchResults{TotalResults: 1, Results: getSearchResults(name)}, nil)
				}
				mtc.EXPECT().SearchPerson("ringo", mock.Anything).Return(nil, errors.New("boom"))
			},
			errFunc: assert.NoError,
		},
		{
			name: "Get Actors By Name - Unhappy No Search Results For Actor",
			args: args{
				called: 3,
				names:  []string{"john", "paul", "george", "ringo"},
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				for _, name := range args.names[:3] {
					mtc.EXPECT().SearchPerson(name, mock.Anything).Return(&go_tmdb.PersonSearchResults{TotalResults: 1, Results: getSearchResults(name)}, nil)
				}
				mtc.EXPECT().SearchPerson("ringo", mock.Anything).Return(&go_tmdb.PersonSearchResults{TotalResults: 0}, nil)
			},
			errFunc: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			m := tmdbMocks.NewMockTMDbClient(t)
			if tt.setup != nil {
				tt.setup(m, tt.args)
			}
			tmdb := tmdb.NewClient(m, logrus.StandardLogger())

			c := make(chan *go_tmdb.Person)
			called := 0
			cb := func() error {
				for range c {
					called += 1
				}
				return nil
			}

			// Act
			err := tmdb.GetActorsByName(context.TODO(), c, cb, tt.args.names...)

			// Assert
			tt.errFunc(t, err)
			assert.Equal(t, tt.args.called, called)
		})
	}

}

func getSearchResults(names ...string) []struct {
	Adult       bool
	ID          int
	Name        string
	Popularity  float32
	ProfilePath string "json:\"profile_path\""
	KnownFor    []struct {
		Adult         bool
		BackdropPath  string "json:\"backdrop_path\""
		ID            int
		OriginalTitle string "json:\"original_title\""
		ReleaseDate   string "json:\"release_date\""
		PosterPath    string "json:\"poster_path\""
		Popularity    float32
		Title         string
		VoteAverage   float32 "json:\"vote_average\""
		VoteCount     uint32  "json:\"vote_count\""
		MediaType     string  "json:\"media_type\""
	} "json:\"known_for\""
} {
	results := make([]struct {
		Adult       bool
		ID          int
		Name        string
		Popularity  float32
		ProfilePath string "json:\"profile_path\""
		KnownFor    []struct {
			Adult         bool
			BackdropPath  string "json:\"backdrop_path\""
			ID            int
			OriginalTitle string "json:\"original_title\""
			ReleaseDate   string "json:\"release_date\""
			PosterPath    string "json:\"poster_path\""
			Popularity    float32
			Title         string
			VoteAverage   float32 "json:\"vote_average\""
			VoteCount     uint32  "json:\"vote_count\""
			MediaType     string  "json:\"media_type\""
		} "json:\"known_for\""
	}, 0, len(names))

	for idx, name := range names {
		results = append(results, struct {
			Adult       bool
			ID          int
			Name        string
			Popularity  float32
			ProfilePath string "json:\"profile_path\""
			KnownFor    []struct {
				Adult         bool
				BackdropPath  string "json:\"backdrop_path\""
				ID            int
				OriginalTitle string "json:\"original_title\""
				ReleaseDate   string "json:\"release_date\""
				PosterPath    string "json:\"poster_path\""
				Popularity    float32
				Title         string
				VoteAverage   float32 "json:\"vote_average\""
				VoteCount     uint32  "json:\"vote_count\""
				MediaType     string  "json:\"media_type\""
			} "json:\"known_for\""
		}{
			ID:   idx,
			Name: name,
		})
	}

	return results
}
