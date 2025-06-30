package tournament

import (
	"discord-bot/src/model"
	"fmt"
	"math/rand"
)

type Storage interface {
	GetAll() ([]model.Member, error)
}

type Tournament struct {
	users   []model.Member
	storage Storage
}

func New(s Storage) Tournament {
	return Tournament{users: []model.Member{}, storage: s}
}

func (t *Tournament) GetPairs() model.LotteryResponse {
	t.users, _ = t.storage.GetAll()

	if len(t.users)%2 != 0 {
		return model.LotteryResponse{Content: "⚠️  No es posible crear un torneo con jugadores impares"}
	}

	rand.Shuffle(len(t.users), func(i, j int) {
		t.users[i], t.users[j] = t.users[j], t.users[i]
	})

	pairs := "🏆  He escogido las llaves para el torneo, y será de la siguiente manera:\n"
	for i := 0; i < len(t.users); i += 2 {
		pairs += fmt.Sprintf("- <@%s> vs <@%s>\n", t.users[i].ID, t.users[i+1].ID)
	}

	return model.LotteryResponse{Content: pairs}
}
