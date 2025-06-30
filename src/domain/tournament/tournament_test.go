package tournament_test

import (
	"discord-bot/src/domain/tournament"
	"discord-bot/src/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAll() ([]model.Member, error) {
	args := m.Called()
	return args.Get(0).([]model.Member), args.Error(1)
}

func TestGetPairs_Even(t *testing.T) {
	mockStorage := new(MockStorage)
	tour := tournament.New(mockStorage)

	members := []model.Member{
		{ID: "1", Name: "user1"},
		{ID: "2", Name: "user2"},
		{ID: "3", Name: "user3"},
		{ID: "4", Name: "user4"},
		{ID: "5", Name: "user5"},
		{ID: "6", Name: "user6"},
		{ID: "7", Name: "user7"},
		{ID: "8", Name: "user8"},
		{ID: "9", Name: "user9"},
		{ID: "10", Name: "user10"},
	}
	mockStorage.On("GetAll").Return(members, nil)

	response := tour.GetPairs()
	assert.Contains(t, response.Content, "vs")
}

func TestGetPairs_Odd(t *testing.T) {
	mockStorage := new(MockStorage)
	tour := tournament.New(mockStorage)

	members := []model.Member{
		{ID: "1", Name: "user1"},
		{ID: "2", Name: "user2"},
		{ID: "3", Name: "user3"},
	}
	mockStorage.On("GetAll").Return(members, nil)

	response := tour.GetPairs()
	assert.Equal(t, "Cannot create pairs with an odd number of players.", response.Content)
}
