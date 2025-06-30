package tournament

import (
	"discord-bot/src/model"
	"discord-bot/src/storage/database"
)

type Storage struct{}

func NewStorage() Storage {
	return Storage{}
}

func (s Storage) GetAll() ([]model.Member, error) {
	var members []model.Member
	result := database.DB().Find(&members)
	return members, result.Error
}
