package lottery

import (
	"discord-bot/src/model"
	"testing"
)

// MockStorage is a mock implementation of the Storage interface for testing purposes.
type MockStorage struct {
	members []model.Member
	winner  model.Winner
}

func (m *MockStorage) GetAll() ([]model.Member, error) {
	return m.members, nil
}

func (m *MockStorage) GetWinner() (model.Winner, error) {
	return m.winner, nil
}

func (m *MockStorage) Add(member *model.Member) error {
	m.members = append(m.members, *member)
	return nil
}

func (m *MockStorage) SetWinner(winner *model.Winner) error {
	m.winner = *winner
	return nil
}

func (m *MockStorage) Clear() error {
	m.members = []model.Member{}
	m.winner = model.Winner{}
	return nil
}

func TestAddUser(t *testing.T) {
	storage := &MockStorage{}
	lottery := New(storage)

	// Test adding a new user
	response := lottery.AddUser("123", "testuser")
	if response.Content != "🎟️  Te has unido al sorteo, suerte testuser!" {
		t.Errorf("Expected '🎟️  Te has unido al sorteo, suerte testuser!', got '%s'", response.Content)
	}

	// Test adding an existing user
	response = lottery.AddUser("123", "testuser")
	if response.Content != "⚠️  Ya te habías inscrito en el sorteo, tramposo" {
		t.Errorf("Expected '⚠️  Ya te habías inscrito en el sorteo, tramposo', got '%s'", response.Content)
	}
}

func TestGetUsers(t *testing.T) {
	storage := &MockStorage{}
	lottery := New(storage)

	// Test with no users
	response := lottery.GetUsers()
	if response.Content != "👀  Todavía no hay participantes! Primero usa `/join`." {
		t.Errorf("Expected '👀  Todavía no hay participantes! Primero usa `/join`.', got '%s'", response.Content)
	}

	// Test with users
	storage.Add(&model.Member{ID: "123", Name: "testuser"})
	response = lottery.GetUsers()
	expectedResponse := "✨ ¡Lista de participantes!\n- <@123>\n"
	if response.Content != expectedResponse {
		t.Errorf("Expected '%s', got '%s'", expectedResponse, response.Content)
	}
}

func TestGetWinner(t *testing.T) {
	storage := &MockStorage{}
	lottery := New(storage)

	// Test with no winner
	response := lottery.GetWinner()
	if response.Content != "🤔  No hay ganador aún, usa `/draw` para iniciar el sorteo." {
		t.Errorf("Expected '🤔  No hay ganador aún, usa `/draw` para iniciar el sorteo.', got '%s'", response.Content)
	}

	// Test with a winner
	storage.SetWinner(&model.Winner{ID: "123", Name: "testuser"})
	response = lottery.GetWinner()
	if response.Content != "🥳  El último ganador fue <@123>" {
		t.Errorf("Expected '🥳  El último ganador fue <@123>', got '%s'", response.Content)
	}
}

func TestDrawWinner(t *testing.T) {
	storage := &MockStorage{}
	lottery := New(storage)

	// Test with no users
	response := lottery.DrawWinner()
	if response.Content != "☹️  Todavía no hay participantes! Primero usa `/join`" {
		t.Errorf("Expected '☹️  Todavía no hay participantes! Primero usa `/join`', got '%s'", response.Content)
	}

	// Test with users
	storage.Add(&model.Member{ID: "123", Name: "testuser"})
	response = lottery.DrawWinner()
	expectedContent := "🎉 ¡Felicidades <@123>! Has ganado el sorteo por 1380RP. 🎊"
	if response.Content != expectedContent {
		t.Errorf("Expected '%s', got '%s'", expectedContent, response.Content)
	}
	if !response.IsAttachment {
		t.Error("Expected IsAttachment to be true")
	}
}

func TestClear(t *testing.T) {
	storage := &MockStorage{}
	lottery := New(storage)

	response := lottery.Clear()
	if response.Content != "🧹  Has limpiado la lista de participantes" {
		t.Errorf("Expected '🧹  Has limpiado la lista de participantes', got '%s'", response.Content)
	}
}
