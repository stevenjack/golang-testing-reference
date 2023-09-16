package repository_test

import (
	"errors"
	"testing"

	"github.com/stevenjack/golang-testing-reference/repository"
	"github.com/stevenjack/golang-testing-reference/repository/domain"
	"github.com/stevenjack/golang-testing-reference/storage"
	"github.com/stevenjack/golang-testing-reference/storage/storagetest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBookFetchByIDWithMock(t *testing.T) {
	t.Parallel()

	errStorage := errors.New("storage error")

	testCases := map[string]struct {
		configureStorageMock func(*storagetest.MockFetcher)
		errAssertion         assert.ErrorAssertionFunc
		expected             domain.Book
		id                   string
	}{
		"When a successful response is returned from storage, return book": {
			configureStorageMock: func(m *storagetest.MockFetcher) {
				m.EXPECT().
					Fetch("1").
					Return([]byte(`{"id":"1","title":"Foo","author":"Bar"}`), nil)
			},
			errAssertion: assert.NoError,
			expected: domain.Book{
				ID:     "1",
				Title:  "Foo",
				Author: "Bar",
			},
			id: "1",
		},
		"When storage errors, return error": {
			configureStorageMock: func(m *storagetest.MockFetcher) {
				m.EXPECT().
					Fetch("1").
					Return(nil, errStorage)
			},
			errAssertion: IsSentinelError(errStorage),
			id:           "1",
		},
		"When storage returns not found error, return specific book not found error": {
			configureStorageMock: func(m *storagetest.MockFetcher) {
				m.EXPECT().
					Fetch("1").
					Return(nil, storage.ErrNotFound)
			},
			errAssertion: IsSentinelError(repository.ErrBookNotFound),
			id:           "1",
		},

		"When unmarshaling response from storage fails, return error": {
			configureStorageMock: func(m *storagetest.MockFetcher) {
				m.EXPECT().
					Fetch("1").
					Return([]byte(`invalid json`), nil)
			},
			errAssertion: assert.Error,
			id:           "1",
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mock := storagetest.NewMockFetcher(ctrl)

			if fn := tc.configureStorageMock; fn != nil {
				fn(mock)
			}

			repo := repository.NewBook(mock)

			book, err := repo.FetchByID(tc.id)
			tc.errAssertion(t, err)
			assert.Equal(t, tc.expected, book)
		})
	}
}

type storageStub struct {
	storage.Fetcher
	FetchFunc func(id string) ([]byte, error)
}

func (m storageStub) Fetch(id string) ([]byte, error) {
	if m.FetchFunc != nil {
		return m.FetchFunc(id)
	}

	return []byte("{}"), nil
}

func TestBookFetchByIDWithoutMock(t *testing.T) {
	t.Parallel()

	errStorage := errors.New("storage error")

	testCases := map[string]struct {
		configureStorageStub func(*storageStub)
		errAssertion         assert.ErrorAssertionFunc
		expected             domain.Book
		id                   string
	}{
		"When a successful response is returned from storage, return book": {
			configureStorageStub: func(m *storageStub) {
				m.FetchFunc = func(id string) ([]byte, error) {
					return []byte(`{"id":"1","title":"Foo","author":"Bar"}`), nil
				}
			},
			errAssertion: assert.NoError,
			expected: domain.Book{
				ID:     "1",
				Title:  "Foo",
				Author: "Bar",
			},
			id: "1",
		},
		"When storage errors, return error": {
			configureStorageStub: func(m *storageStub) {
				m.FetchFunc = func(id string) ([]byte, error) {
					return nil, errStorage
				}
			},
			errAssertion: IsSentinelError(errStorage),
			id:           "1",
		},
		"When storage returns not found error, return specific book not found error": {
			configureStorageStub: func(m *storageStub) {
				m.FetchFunc = func(id string) ([]byte, error) {
					return nil, storage.ErrNotFound
				}
			},
			errAssertion: IsSentinelError(repository.ErrBookNotFound),
			id:           "1",
		},
		"When unmarshaling response from storage fails, return error": {
			configureStorageStub: func(m *storageStub) {
				m.FetchFunc = func(id string) ([]byte, error) {
					return []byte(`invalid json`), nil
				}
			},
			errAssertion: assert.Error,
			id:           "1",
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			stub := storageStub{}

			if fn := tc.configureStorageStub; fn != nil {
				fn(&stub)
			}

			repo := repository.NewBook(stub)

			book, err := repo.FetchByID(tc.id)
			tc.errAssertion(t, err)
			assert.Equal(t, tc.expected, book)
		})
	}
}

func IsSentinelError(expectedErr error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, params ...interface{}) bool {
		return assert.ErrorIs(t, err, expectedErr, params...)
	}
}
