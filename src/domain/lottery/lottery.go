package lottery

import (
	"discord-bot/src/model"
	"fmt"
	"math/rand/v2"
)

type Storage interface {
	GetAll() ([]model.Member, error)
	GetWinner() (model.Winner, error)
	Add(member *model.Member) error
	SetWinner(winner *model.Winner) error
	Clear() error
}

type Lottery struct {
	users   []model.Member
	winner  model.Winner
	storage Storage
}

func New(s Storage) Lottery {
	return Lottery{users: []model.Member{}, winner: model.Winner{}, storage: s}
}

func (l *Lottery) AddUser(userId string, userName string) model.LotteryResponse {
	l.users, _ = l.storage.GetAll()
	member := model.Member{Name: userName, ID: userId}

	for _, user := range l.users {
		if user.ID == member.ID {
			return model.LotteryResponse{Content: "⚠️  Ya te habías inscrito en el sorteo, tramposo"}
		}
	}

	l.users = append(l.users, member)
	l.storage.Add(&member)

	return model.LotteryResponse{Content: fmt.Sprintf("🎟️  Te has unido al sorteo, suerte %s!", userName)}
}

func (l *Lottery) GetUsers() model.LotteryResponse {
	l.users, _ = l.storage.GetAll()

	if len(l.users) == 0 {
		return model.LotteryResponse{Content: "👀  Todavía no hay participantes! Primero usa `/join`."}
	}

	result := "✨ ¡Lista de participantes!\n"
	for _, user := range l.users {
		result += "- <@" + user.ID + ">\n"
	}

	return model.LotteryResponse{Content: result}
}

func (l *Lottery) GetWinner() model.LotteryResponse {
	l.winner, _ = l.storage.GetWinner()

	if l.winner.ID == "" {
		return model.LotteryResponse{Content: "🤔  No hay ganador aún, usa `/draw` para iniciar el sorteo."}
	}

	return model.LotteryResponse{Content: fmt.Sprintf("🥳  El último ganador fue <@%s>", l.winner.ID)}
}

func (l *Lottery) DrawWinner() model.LotteryResponse {
	l.users, _ = l.storage.GetAll()

	if len(l.users) == 0 {
		return model.LotteryResponse{Content: "☹️  Todavía no hay participantes! Primero usa `/join`"}
	}

	randomIndex := rand.IntN(len(l.users))

	l.winner.ID = l.users[randomIndex].ID
	l.winner.Name = l.users[randomIndex].Name

	l.users = []model.Member{}
	l.storage.SetWinner(&l.winner)
	l.storage.Clear()

	return model.LotteryResponse{
		IsAttachment: true,
		Content:      fmt.Sprintf("🎉 ¡Felicidades <@%s>! Has ganado el sorteo por 1380RP. 🎊", l.winner.ID),
	}
}
