# chain-react

chain-react is a web based trivia game for guessing a path between actors who have acted in common movies.

[![CI](https://github.com/jj-style/chain-react/actions/workflows/ci.yml/badge.svg)](https://github.com/jj-style/chain-react/actions/workflows/ci.yml)

## Demo
You can access the site here https://chain-react.xzy.

## Development

### Requirements
- [`go`](https://go.dev/)
- [`react`](https://react.dev/)
- [`docker`](https://docs.docker.com/get-docker/)
- create a [TMDB api key](https://developer.themoviedb.org/docs)

### Running the server
- Copy the [.chain-react.yaml.example](backend/.chain-react.yaml.example) file to `backend` and fill in the required fields.
- `go run main.go server`

### Running the frontend
- Copy the [.env.local.example](frontend/.env.local.example) to `frontend` and fill in the variables.
- `npm install && npm run start:local`

### Start other services
- Copy the [.env.example](docker/.env.example) file to the `docker` folder
- `docker compose up -d`
  - you may want to comment out `proxy` and `backend` services while developing locally
