package user

import (
	"discord-bot/src/model"
	"discord-bot/src/storage/database"
)

type Storage interface {
	GetAll() ([]model.Member, error)
	GetWinner() (model.Winner, error)
	Add(member *model.Member) error
	SetWinner(member *model.Winner) error
	Clear() error
}

type Store struct{}

func NewStorage() *Store {
	return &Store{}
}

func (s *Store) GetAll() ([]model.Member, error) {
	members := make([]model.Member, 0)
	database.DB().Find(&members)

	return members, nil
}

func (s *Store) GetWinner() (model.Winner, error) {
	winner := model.Winner{}
	database.DB().Find(&winner)

	return winner, nil
}

func (s *Store) Add(member *model.Member) error {
	if result := database.DB().Create(&member); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store) SetWinner(winner *model.Winner) error {
	database.DB().Unscoped().Where("1 = 1").Delete(&model.Winner{})

	if result := database.DB().Create(&winner); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store) Clear() error {
	database.DB().Unscoped().Where("1 = 1").Delete(&model.Member{})
	return nil
}
