package search

//go:generate mockery
type Repository interface {
	AddDocuments(docs interface{}, index string) error
}
