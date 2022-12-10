package main

import (
	"github.com/DenChenn/yenno-dc/config"
	"github.com/DenChenn/yenno-dc/dao"
	"github.com/DenChenn/yenno-dc/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
	"github.com/teris-io/shortid"
	"log"
	"os"
	"os/signal"
)

func init() {
	config.Env = config.Load()
}

func main() {
	s, err := discordgo.New("Bot " + config.Env.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	if err := s.Open(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// create new handler
	newDAO := dao.New(redis.NewClient(&redis.Options{
		Addr:     config.Env.RedisUrl,
		Password: config.Env.RedisPassword,
	}))
	generator, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		log.Fatal(err)
	}
	newHandler := handler.NewHandler(newDAO, generator)

	// Register commands
	commands := registerCommand(s)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		newHandler.Select(i)(s, i)
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	// Gracefully shut down
	log.Println("Removing commands...")
	removeCommand(s, commands)
	log.Println("Gracefully shutting down.")
}

func registerCommand(s *discordgo.Session) []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-deployment-config",
			Description: "Create config for deployment",
		},
		{
			Name:        "list-all-deployment-config",
			Description: "List all deployment config",
		},
		{
			Name:        "delete-deployment-config",
			Description: "Delete deployment config",
		},
		{
			Name:        "deploy-with-config",
			Description: "Deploy with config",
		},
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, c := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Guilds[0].ID, c)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", c.Name, err)
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands
}

func removeCommand(s *discordgo.Session, registeredCommand []*discordgo.ApplicationCommand) {
	for _, c := range registeredCommand {
		if err := s.ApplicationCommandDelete(s.State.User.ID, s.State.Guilds[0].ID, c.ID); err != nil {
			log.Panicf("Cannot delete '%v' command: %v", c.Name, err)
		}
	}
}
