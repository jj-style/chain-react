package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jj-style/chain-react/src/config"
	"github.com/jj-style/chain-react/src/db"
	"github.com/jj-style/chain-react/src/search"
	"github.com/jj-style/chain-react/src/tmdb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func givenServer(mockDb *db.MockRepository, mockTMDb *tmdb.MockTMDb, mockSearch *search.MockRepository) Server {
	// setup server
	srv := Server{
		Router: newRouter(true),
		Config: &config.RConfig{
			Repo:   mockDb,
			Tmdb:   mockTMDb,
			Search: mockSearch,
		},
		Log: logrus.New(),
	}
	srv.setupRoutes()
	return srv
}

func TestGetRandomActor(t *testing.T) {
	// setup mocks
	mockDb := db.NewMockRepository(t)
	mockTMDb := tmdb.NewMockTMDb(t)
	mockSearch := search.NewMockRepository(t)

	actor := db.Actor{Id: 1, Name: "Jackie Chan", Popularity: 100}
	mockDb.EXPECT().RandomActor().Return(&actor, nil)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/randomActor", nil)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)

	var got db.Actor
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, actor, got)
}

func TestGetRandomActorNot(t *testing.T) {
	// setup mocks
	mockDb := db.NewMockRepository(t)
	mockTMDb := tmdb.NewMockTMDb(t)
	mockSearch := search.NewMockRepository(t)

	jc := db.Actor{Id: 1, Name: "Jackie Chan", Popularity: 100}
	ct := db.Actor{Id: 2, Name: "Chris Tucker", Popularity: 100}
	mockDb.EXPECT().RandomActorNotId(1).Return(&ct, nil)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/randomActorNot/%d", jc.Id), nil)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)

	var got db.Actor
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, ct, got)
}
