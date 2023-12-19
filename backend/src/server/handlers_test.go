package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jj-style/chain-react/src/db"
	dbMocks "github.com/jj-style/chain-react/src/db/mocks"
	searchMocks "github.com/jj-style/chain-react/src/search/mocks"
	tmdbMocks "github.com/jj-style/chain-react/src/tmdb/mocks"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func givenServer(mockDb *dbMocks.MockRepository, mockTMDb *tmdbMocks.MockTMDb, mockSearch *searchMocks.MockRepository) Server {
	// setup server
	srv := Server{
		Router: newRouter(true, "*"),
		Repo:   mockDb,
		Tmdb:   mockTMDb,
		Search: mockSearch,
		Log:    logrus.New(),
	}
	srv.setupRoutes()
	return srv
}

func TestGetRandomActor_Happy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

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

func TestGetRandomActor_Error(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().RandomActor().Return(nil, errors.New("boom"))

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/randomActor", nil)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetRandomActorNot_Happy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

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

func TestGetRandomActorNot_InvalidIdParam(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/randomActorNot/notAnInt", nil)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRandomActorNot_Error(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().RandomActorNotId(1).Return(nil, errors.New("boom"))

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/randomActorNot/1", nil)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyEdgesHappy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	jc := db.Actor{Id: 1, Name: "Jackie Chan", Popularity: 100}
	ct := db.Actor{Id: 2, Name: "Chris Tucker", Popularity: 100}
	mov := db.Movie{Id: 1, Title: "Rush Hour"}
	mockDb.EXPECT().VerifyWithEdges(db.Chain{1, 2}).Return(
		[]*db.Edge{
			{
				Src: db.Credit{
					Actor:     jc,
					Movie:     mov,
					Character: "Chief Inspector Lee",
				},
				Dest: db.Credit{
					Actor:     ct,
					Movie:     mov,
					Character: "James Carter",
				},
			},
		},
		nil,
	)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2]}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verifyEdges", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)

	type response struct {
		Valid bool       `json:"valid"`
		Error string     `json:"error"`
		Chain []*db.Edge `json:"chain"`
	}
	var got response
	err := json.Unmarshal(w.Body.Bytes(), &got)

	assert.NoError(t, err)
	assert.Equal(t, true, got.Valid)
	assert.Equal(t, "", got.Error)
	assert.Equal(t, 1, len(got.Chain))
}

func TestVerifyEdgesUnhappy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().VerifyWithEdges(db.Chain{1, 2}).Return(
		[]*db.Edge{},
		errors.New("boom"),
	)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2]}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verifyEdges", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	type response struct {
		Valid bool       `json:"valid"`
		Error string     `json:"error"`
		Chain []*db.Edge `json:"chain"`
	}
	var got response
	err := json.Unmarshal(w.Body.Bytes(), &got)

	assert.NoError(t, err)
	assert.Equal(t, false, got.Valid)
	assert.Equal(t, "boom", got.Error)
}

func TestVerifyEdgesBadRequest(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": "not slice of ints"}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verifyEdges", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyHappy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().Verify(db.Chain{1, 2}).Return(true, nil)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2]}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verify", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)

	type response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}
	var got response
	err := json.Unmarshal(w.Body.Bytes(), &got)

	assert.NoError(t, err)
	assert.Equal(t, true, got.Valid)
	assert.Equal(t, "", got.Error)
}

func TestVerifyUnhappy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().Verify(db.Chain{1, 2}).Return(false, errors.New("boom"))

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2]}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verify", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	type response struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}
	var got response
	err := json.Unmarshal(w.Body.Bytes(), &got)

	assert.NoError(t, err)
	assert.Equal(t, false, got.Valid)
	assert.Equal(t, "boom", got.Error)
}

func TestVerifyBadRequest(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": "not slice of ints"}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/verify", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetGraph_UnhappyBadRequestJson(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"bad": "json"}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/graph", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetGraph_UnhappyBadRequestLength(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2], "length": 10}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/graph", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetGraph_UnhappyError(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().GetGraph(4, 1, 2).Return([]dbtype.Path{}, errors.New("boom"))

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2], "length": 4}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/graph", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetGraph_Happy(t *testing.T) {
	// setup mocks
	mockDb := dbMocks.NewMockRepository(t)
	mockTMDb := tmdbMocks.NewMockTMDb(t)
	mockSearch := searchMocks.NewMockRepository(t)

	mockDb.EXPECT().GetGraph(4, 1, 2).Return([]dbtype.Path{
		{
			Nodes: []dbtype.Node{
				{ElementId: "n1"},
				{ElementId: "n2"},
			},
			Relationships: []dbtype.Relationship{
				{
					StartElementId: "n1",
					EndElementId:   "n2",
					ElementId:      "r1",
					Type:           "relationshipType",
				},
			},
		},
	}, nil)

	// create server
	srv := givenServer(mockDb, mockTMDb, mockSearch)

	// handle request
	jsonBody := []byte(`{"chain": [1, 2], "length": 4}`)
	bodyReader := bytes.NewReader(jsonBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/graph", bodyReader)
	srv.Router.ServeHTTP(w, req)

	// check response
	assert.Equal(t, http.StatusOK, w.Code)
}
