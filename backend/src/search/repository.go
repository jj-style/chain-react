package search

type Repository interface {
	AddDocuments(docs interface{}, index string) error
}
