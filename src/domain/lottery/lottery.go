package lottery

import (
	"discord-bot/src/model"
	"fmt"
	"math/rand/v2"
)

type Lottery struct {
	users  []model.Member
	winner model.Member
}

func New() Lottery {
	return Lottery{users: []model.Member{}, winner: model.Member{}}
}

func (l *Lottery) AddUser(userId string, userName string) model.LotteryResponse {
	for _, user := range l.users {
		if user.ID == userId {
			return model.LotteryResponse{Content: "⚠️  Ya te habías inscrito en el sorteo, tramposo"}
		}
	}

	l.users = append(l.users, model.Member{Name: userName, ID: userId})

	return model.LotteryResponse{Content: fmt.Sprintf("🎟️  Te has unido al sorteo, suerte %s!", userName)}
}

func (l *Lottery) GetUsers() model.LotteryResponse {
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
	if l.winner.ID == "" {
		return model.LotteryResponse{Content: "🤔  No hay ganador aún, usa `/draw` para iniciar el sorteo."}
	}

	return model.LotteryResponse{Content: fmt.Sprintf("🥳  El último ganador fue <@%s>", l.winner.ID)}
}

func (l *Lottery) DrawWinner() model.LotteryResponse {
	if len(l.users) == 0 {
		return model.LotteryResponse{Content: "☹️  Todavía no hay participantes! Primero usa `/join`"}
	}
	randomIndex := rand.IntN(len(l.users))
	l.winner = l.users[randomIndex]
	l.users = []model.Member{}

	return model.LotteryResponse{
		IsAttachment: true,
		Content:      fmt.Sprintf("🎉 ¡Felicidades <@%s>! Has ganado el sorteo por 1380RP. 🎊", l.winner.ID),
	}
}
