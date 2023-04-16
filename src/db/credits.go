package db

type CreditIn struct {
	ActorId   int
	MovieId   int
	CreditId  string
	Character string
}

type Credit struct {
	Actor
	Movie
	CreditId  string
	Character string
}
