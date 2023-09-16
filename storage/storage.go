package storage

//go:generate go run go.uber.org/mock/mockgen -package storagetest -destination=storagetest/mocks.go -source=storage.go
type Fetcher interface {
	Fetch(id string) ([]byte, error)
}
