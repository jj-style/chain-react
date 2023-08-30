package tmdb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jj-style/chain-react/src/tmdb"
	tmdbMocks "github.com/jj-style/chain-react/src/tmdb/mocks"
	go_tmdb "github.com/jj-style/go-tmdb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetActorsMovieCredits(t *testing.T) {
	type args struct {
		id         int
		credits    *go_tmdb.PersonMovieCredits
		creditsErr error
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*tmdbMocks.MockTMDbClient, args)
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "Get Actors Movie Credits - Happy",
			args: args{
				id: 1,
				credits: &go_tmdb.PersonMovieCredits{
					ID:   1,
					Cast: getCast(map[string]string{"Rush Hour": "Lee", "Kung Fu Panda": "Monkey"}),
				},
				creditsErr: nil,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonMovieCredits(args.id, mock.Anything).Return(args.credits, args.creditsErr)
			},
			errFunc: assert.NoError,
		},
		{
			name: "Get Actors Movie Credits - Unhappy error",
			args: args{
				id:         1,
				credits:    nil,
				creditsErr: errors.New("boom"),
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonMovieCredits(args.id, mock.Anything).Return(args.credits, args.creditsErr)
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

			// Act
			_, err := tmdb.GetActorMovieCredits(context.TODO(), tt.args.id)

			// Assert
			tt.errFunc(t, err)
		})
	}

}

func TestGetAllActorMovieCredits(t *testing.T) {
	type args struct {
		ids    []int
		called int
	}
	tests := []struct {
		name    string
		args    args
		setup   func(*tmdbMocks.MockTMDbClient, args)
		errFunc assert.ErrorAssertionFunc
	}{
		{
			name: "Get All Actor Movie Credits - Happy",
			args: args{
				ids:    []int{1, 2, 3},
				called: 3,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonMovieCredits(1, mock.Anything).Return(&go_tmdb.PersonMovieCredits{
					Cast: getCast(map[string]string{"Rush Hour": "Lee", "Kung Fu Panda": "Monkey"}),
				}, nil).Once()
				mtc.EXPECT().GetPersonMovieCredits(2, mock.Anything).Return(&go_tmdb.PersonMovieCredits{
					Cast: getCast(map[string]string{"The Departed": "Colin", "The Bourne Identity": "Jason Bourne"}),
				}, nil).Once()
				mtc.EXPECT().GetPersonMovieCredits(3, mock.Anything).Return(&go_tmdb.PersonMovieCredits{
					Cast: getCast(map[string]string{"The Departed": "Billy", "Inception ": "Cobb"}),
				}, nil).Once()
			},
			errFunc: assert.NoError,
		},
		{
			name: "Get All Actor Movie Credits - Unhappy Error Getting Credits Of One Actor",
			args: args{
				ids:    []int{1, 2, 3},
				called: 2,
			},
			setup: func(mtc *tmdbMocks.MockTMDbClient, args args) {
				mtc.EXPECT().GetPersonMovieCredits(1, mock.Anything).Return(&go_tmdb.PersonMovieCredits{
					Cast: getCast(map[string]string{"Rush Hour": "Lee", "Kung Fu Panda": "Monkey"}),
				}, nil).Once()
				mtc.EXPECT().GetPersonMovieCredits(2, mock.Anything).Return(&go_tmdb.PersonMovieCredits{
					Cast: getCast(map[string]string{"The Departed": "Colin", "The Bourne Identity": "Jason Bourne"}),
				}, nil).Once()
				mtc.EXPECT().GetPersonMovieCredits(3, mock.Anything).Return(nil, errors.New("boom")).Once()
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

			c := make(chan *go_tmdb.PersonMovieCredits)
			called := 0
			cb := func() error {
				for range c {
					called += 1
				}
				return nil
			}

			// Act
			err := tmdb.GetAllActorMovieCredits(context.TODO(), c, cb, tt.args.ids...)

			// Assert
			tt.errFunc(t, err)
			assert.Equal(t, tt.args.called, called)
		})
	}

}

func getCast(credits map[string]string) []struct {
	Adult         bool
	Character     string
	CreditID      string "json:\"credit_id\""
	ID            int
	OriginalTitle string "json:\"original_title\""
	PosterPath    string "json:\"poster_path\""
	ReleaseDate   string "json:\"release_date\""
	Title         string
} {
	cast := make([]struct {
		Adult         bool
		Character     string
		CreditID      string "json:\"credit_id\""
		ID            int
		OriginalTitle string "json:\"original_title\""
		PosterPath    string "json:\"poster_path\""
		ReleaseDate   string "json:\"release_date\""
		Title         string
	}, 0, len(credits))
	id := 1
	for title, character := range credits {
		cast = append(cast, struct {
			Adult         bool
			Character     string
			CreditID      string "json:\"credit_id\""
			ID            int
			OriginalTitle string "json:\"original_title\""
			PosterPath    string "json:\"poster_path\""
			ReleaseDate   string "json:\"release_date\""
			Title         string
		}{
			ID:        id,
			Title:     title,
			Character: character,
		})
		id += 1
	}
	return cast
}
