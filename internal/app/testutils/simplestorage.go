package testutils

import (
	"github.com/nayakunin/shortener/internal/app/interfaces"
)

type GetMock struct {
	Success string
	Error   error
}

type AddMock struct {
	Success string
	Error   error
}

type DeleteUserUrlsMock struct {
	Error error
}

type GetUrlsByUserMock struct {
	Success map[string]string
	Error   error
}

type AddBatchMock struct {
	Success []interfaces.DBBatchOutput
	Error   error
}

type StatsMock struct {
	Success interfaces.Stats
	Error   error
}

type SimpleMockStorageParameters struct {
	Get            GetMock
	Add            AddMock
	AddBatch       AddBatchMock
	DeleteUserUrls DeleteUserUrlsMock
	GetUrlsByUser  GetUrlsByUserMock
	Stats          StatsMock
}

// SimpleMockStorage is a mock struct for interfaces.Storage
type SimpleMockStorage struct {
	parameters SimpleMockStorageParameters
}

// NewSimpleMockStorage creates a new mock storage
func NewSimpleMockStorage(parameters SimpleMockStorageParameters) interfaces.Storage {
	return &SimpleMockStorage{
		parameters: parameters,
	}
}

// Get implements storage.Storager
func (s *SimpleMockStorage) Get(key string) (string, error) {
	return s.parameters.Get.Success, s.parameters.Get.Error
}

// Add implements storage.Storager
func (s *SimpleMockStorage) Add(link string, userID string) (string, error) {
	return s.parameters.Add.Success, s.parameters.Add.Error
}

// GetUrlsByUser implements interfaces.Storage
func (s *SimpleMockStorage) GetUrlsByUser(userID string) (map[string]string, error) {
	return s.parameters.GetUrlsByUser.Success, s.parameters.GetUrlsByUser.Error
}

// AddBatch implements interfaces.Storage
func (s *SimpleMockStorage) AddBatch(batches []interfaces.BatchInput, userID string) ([]interfaces.DBBatchOutput, error) {
	return s.parameters.AddBatch.Success, s.parameters.AddBatch.Error
}

// DeleteUserUrls implements interfaces.Storage
func (s *SimpleMockStorage) DeleteUserUrls(userID string, keys []string) error {
	return s.parameters.DeleteUserUrls.Error
}

func (s *SimpleMockStorage) Stats() (interfaces.Stats, error) {
	return s.parameters.Stats.Success, s.parameters.Stats.Error
}
