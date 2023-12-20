package gamemanager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	db2 "github.com/jj-style/chain-react/src/db"
	db "github.com/jj-style/chain-react/src/db/mocks"
	gamemanager "github.com/jj-style/chain-react/src/server/game_manager"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name     string
		args     args
		setup    func(*db.MockRepository, redismock.ClientMock)
		checkErr assert.ErrorAssertionFunc
		want     *gamemanager.Game
	}{
		{
			name: "happy path",
			setup: func(repo *db.MockRepository, cache redismock.ClientMock) {
				repo.EXPECT().RandomActor().Return(&db2.Actor{Id: 1, Name: "a"}, nil)
				repo.EXPECT().RandomActorNotId(1).Return(&db2.Actor{Id: 2, Name: "b"}, nil)

				value := []byte(`{"start":{"id":1,"name":"a","popularity":0},"end":{"id":2,"name":"b","popularity":0}}`)

				cache.ExpectGet(gamemanager.GameKey).SetVal("different game")
				cache.ExpectSet(gamemanager.GameKey, value, 0).RedisNil()
			},
			checkErr: assert.NoError,
			want:     &gamemanager.Game{Start: &db2.Actor{Id: 1, Name: "a", Popularity: 0}, End: &db2.Actor{Id: 2, Name: "b", Popularity: 0}},
		},
		{
			name: "error getting start",
			setup: func(repo *db.MockRepository, cache redismock.ClientMock) {
				repo.EXPECT().RandomActor().Return(nil, errors.New("boom"))
			},
			checkErr: assert.Error,
		},
		{
			name: "error getting end",
			setup: func(repo *db.MockRepository, cache redismock.ClientMock) {
				repo.EXPECT().RandomActor().Return(&db2.Actor{Id: 1, Name: "a"}, nil)
				repo.EXPECT().RandomActorNotId(1).Return(nil, errors.New("boom"))
			},
			checkErr: assert.Error,
		},
		{
			name: "error getting current game from redis",
			setup: func(repo *db.MockRepository, cache redismock.ClientMock) {
				repo.EXPECT().RandomActor().Return(&db2.Actor{Id: 1, Name: "a"}, nil)
				repo.EXPECT().RandomActorNotId(1).Return(&db2.Actor{Id: 2, Name: "b"}, nil)
				cache.ExpectGet(gamemanager.GameKey).SetErr(errors.New("boom"))
			},
			checkErr: assert.Error,
		},
		{
			name: "error setting new game in redis",
			setup: func(repo *db.MockRepository, cache redismock.ClientMock) {
				repo.EXPECT().RandomActor().Return(&db2.Actor{Id: 1, Name: "a"}, nil)
				repo.EXPECT().RandomActorNotId(1).Return(&db2.Actor{Id: 2, Name: "b"}, nil)
				value := []byte(`{"start":{"id":1,"name":"a","popularity":0},"end":{"id":2,"name":"b","popularity":0}}`)
				cache.ExpectGet(gamemanager.GameKey).SetVal("different game")
				cache.ExpectSet(gamemanager.GameKey, value, 0).SetErr(errors.New("boom"))
			},
			checkErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cache, redisMock := redismock.NewClientMock()
			repo := db.NewMockRepository(t)
			logger := log.New()

			if tt.setup != nil {
				tt.setup(repo, redisMock)
			}

			got, err := gamemanager.NewGame(context.Background(), logger, repo, cache)
			tt.checkErr(t, err)
			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}
			if err := redisMock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
