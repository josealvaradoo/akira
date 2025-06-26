package discord

import (
	"bytes"
	"discord-bot/src/model"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{}
	handlers = []model.Handler{}
)

type Discord struct {
	token   string
	session *discordgo.Session
}

// Create a new Discord instance with a provided token
func New(token string) Discord {
	if token == "" {
		log.Fatal("a token is required")
	}

	// Create a new Discord session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	return Discord{
		token:   token,
		session: session,
	}
}

func (d *Discord) SetCommands(c []model.Command) *Discord {
	commands = make([]*discordgo.ApplicationCommand, len(c))

	fmt.Println("⚙️ Setting up commands...")

	for i, cmd := range c {
		commands[i] = &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		}
	}

	return d
}

func (d *Discord) AddHandler(h model.Handler) *Discord {
	if len(commands) == 0 {
		panic("⚠  commands must be set before adding handlers")
	}

	handlers = append(handlers, h)
	return d
}

// Create a new Discord session and start the bot
func (d *Discord) Start() *Discord {
	d.session.AddHandler(ready)
	d.session.AddHandler(interactionCreate)

	handleWebsocket(d.session)

	return d
}

func handleWebsocket(session *discordgo.Session) {
	if err := session.Open(); err != nil {
		log.Fatal("⚠  Error opening connection:", err)
	}

	defer session.Close()

	fmt.Println("✨ Bot is now running. Press CTRL+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("")
	fmt.Println("👋🏻 Gracefully shutting down...")
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("🔐 Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)

	err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Name: "sorteos de RP 🎲",
				Type: discordgo.ActivityTypeGame,
			},
		},
		AFK: false,
	})
	if err != nil {
		log.Printf("Error setting online status: %v", err)
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("❌ Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, handler := range handlers {
		if i.ApplicationCommandData().Name == handler.Name {
			if handler.IsAdminOnly && !isUserAdmin(s, i) {
				respond(s, i, "❌ You need administrator permissions to use this command.", false, false)
				return
			}

			event := handler.Event(i)
			isAttachment := handler.IsAttachment && event.IsAttachment

			fmt.Printf("📥 Received command: /%s (from user: %s)\n", handler.Name, i.Member.User.DisplayName())

			respond(s, i, event.Content, handler.ForEveryone, isAttachment)
		}
	}
}

func respond(session *discordgo.Session, i *discordgo.InteractionCreate, content string, forEveryone bool, isAttachment bool) {
	data := &discordgo.InteractionResponseData{
		Content: content,
	}

	if !forEveryone {
		data.Flags = discordgo.MessageFlagsEphemeral
	}

	if isAttachment {
		data.Files = handleAttachment()
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func handleAttachment() []*discordgo.File {
	fileData, err := os.ReadFile("src/assets/winner.png")
	if err != nil {
		log.Printf("Error reading image file: %v", err)
		return []*discordgo.File{}
	}

	fmt.Println("📎 Sending attachment...")

	return []*discordgo.File{
		{
			Name:        "winner.png",
			ContentType: "image/png",
			Reader:      bytes.NewReader(fileData),
		},
	}
}

func isUserAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	// Check if interaction is from a guild (server)
	if i.GuildID == "" {
		return false
	}

	// Get user's permissions in the guild
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		log.Printf("Error getting user permissions: %v", err)
		return false
	}

	// Check if user has administrator permission
	return permissions&discordgo.PermissionAdministrator != 0
}
