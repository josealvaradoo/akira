package main

import (
	"discord-bot/src/domain/discord"
	"discord-bot/src/domain/lottery"
	"discord-bot/src/model"
	"discord-bot/src/storage/database"
	db "discord-bot/src/storage/database/lottery"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	database.New(database.SQLite)

	database.DB().AutoMigrate(model.Member{}, model.Winner{})

	s := db.NewStorage()
	l := lottery.New(s)
	d := discord.New(os.Getenv("DISCORD_BOT_TOKEN"))

	commands := []model.Command{
		{Name: "rules", Description: "Reglas de Akira", ForEveryone: true, IsAttachment: true, Event: rules},
		{Name: "join", Description: "Únete al sorteo por 1380RP", ForEveryone: false, Event: addUserDecorator(l.AddUser)},
		{Name: "list", Description: "Lista de participantes que se han unido al sorteo", ForEveryone: true, Event: genericDecorator(l.GetUsers)},
		{Name: "winner", Description: "Conoce quién ganó el último sorteo", ForEveryone: false, Event: genericDecorator(l.GetWinner)},
		{Name: "draw", Description: "Iniciar el sorteo", ForEveryone: true, IsAttachment: true, IsAdminOnly: true, Event: genericDecorator(l.DrawWinner)},
	}

	d.SetCommands(commands)

	for _, command := range commands {
		d.AddHandler(model.Handler{
			Name:         command.Name,
			Event:        command.Event,
			ForEveryone:  command.ForEveryone,
			IsAttachment: command.IsAttachment,
			IsAdminOnly:  command.IsAdminOnly,
		})
	}

	d.Start()
}

func rules(u *discordgo.InteractionCreate) model.LotteryResponse {
	content := "Akira es nuestro bot propio para sortear RP en el servidor. Y sus reglas son muy sencillas.\n\n"
	content += "- El premio será el indicado por Heimeñinger\n"
	content += "- El premio es intransferible, y por ahora solo será en League of Legends\n"
	content += "- Solo participarán los miembros del server, que se hayan unido al soteo en el momento del anuncio y aparezcan en `/list`\n"
	content += "- Solo puedes unirte una vez, haciendo uso de `/join`\n\n"
	content += "Los sorteos serán espontáneos, dependen de patrocinio."

	return model.LotteryResponse{Content: content}
}

func addUserDecorator(f func(string, string) model.LotteryResponse) func(*discordgo.InteractionCreate) model.LotteryResponse {
	return func(i *discordgo.InteractionCreate) model.LotteryResponse {
		return f(i.Member.User.ID, i.Member.User.DisplayName())
	}
}

func genericDecorator(f func() model.LotteryResponse) func(*discordgo.InteractionCreate) model.LotteryResponse {
	return func(i *discordgo.InteractionCreate) model.LotteryResponse {
		return f()
	}
}
