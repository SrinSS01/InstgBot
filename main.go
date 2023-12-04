package main

import (
	"InstgBot/commands"
	"InstgBot/config"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

type Execute struct {
	Slash func(*discordgo.Session, *discordgo.InteractionCreate)
	Dash  func(*discordgo.Session, *discordgo.MessageCreate, string)
}

var (
	discord *discordgo.Session
	cnfg    = config.Config{}
	//accRegex      = regexp.MustCompile("^-acc(?: +(?P<sessionid>.+))?$")
	dashCommandRegex = regexp.MustCompile("-(?P<name>[\\w-]+)(?:\\s+(?P<args>.+))?")
	cmds             = []*discordgo.ApplicationCommand{
		commands.Acc.Command,
		commands.CancelAcc.Command,
		commands.AddSessionId.Command,
		commands.Copy.Command,
	}
	commandHandlers = map[string]*Execute{
		commands.Acc.Command.Name: {
			Slash: commands.Acc.ExecuteSlash,
			Dash:  commands.Acc.ExecuteDash,
		},
		commands.CancelAcc.Command.Name: {
			Slash: commands.CancelAcc.ExecuteSlash,
			Dash:  commands.CancelAcc.ExecuteDash,
		},
		commands.AddSessionId.Command.Name: {
			Slash: commands.AddSessionId.ExecuteSlash,
		},
		commands.Copy.Command.Name: {
			Slash: commands.Copy.ExecuteSlash,
			Dash:  commands.Copy.ExecuteDash,
		},
	}
)

func init() {
	file, err := os.ReadFile("session.json")
	if err != nil {
		saveSession()
		return
	}
	if err := json.Unmarshal(file, &commands.SessionObj); err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Print("Enter bot token: ")
		if _, err := fmt.Scanln(&cnfg.Token); err != nil {
			log.Fatal("Error during Scanln(): ", err)
		}
		print("Enter api key: ")
		if _, err := fmt.Scanln(&cnfg.ApiKey); err != nil {
			log.Fatal("Error during Scanln(): ", err.Error())
		}
		print("Enter api url: ")
		if _, err := fmt.Scanln(&cnfg.ApiUrl); err != nil {
			log.Fatal("Error during Scanln(): ", err.Error())
		}
		configJson()
		return
	}
	if err := json.Unmarshal(file, &cnfg); err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func configJson() {
	marshal, err := json.Marshal(&cnfg)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
		return
	}
	if err := os.WriteFile("config.json", marshal, 0644); err != nil {
		log.Fatal("Error during WriteFile(): ", err)
	}
}

func saveSession() {
	marshal, err := json.Marshal(&commands.SessionObj)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
		return
	}
	if err := os.WriteFile("session.json", marshal, 0644); err != nil {
		log.Fatal("Error during WriteFile(): ", err)
	}
	return
}

func onReady(s *discordgo.Session, _ *discordgo.Ready) {
	log.Println(s.State.User.Username + " is ready")
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	matches := dashCommandRegex.FindStringSubmatch(m.Content)
	if len(matches) == 0 {
		return
	}
	name := matches[dashCommandRegex.SubexpIndex("name")]
	args := matches[dashCommandRegex.SubexpIndex("args")]
	f := commandHandlers[name]
	if f == nil {
		return
	}
	_, _ = s.ChannelMessageSendReply(m.ChannelID, "Processing...", m.Reference())
	f.Dash(s, m, args)
}

func slashCommandInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	commandHandlers[interaction.ApplicationCommandData().Name].Slash(session, interaction)
}

func main() {
	commands.Acc.Config = &cnfg
	var err error
	discord, err = discordgo.New("Bot " + cnfg.Token)
	if err != nil {
		log.Fatal("Error creating Discord session", err)
		return
	}
	discord.Identify.Intents |= discordgo.IntentMessageContent
	defer discord.AddHandler(onReady)()
	defer discord.AddHandler(messageCreate)()
	defer discord.AddHandler(commands.Acc.QuestionAnswerHandler)()
	defer discord.AddHandler(slashCommandInteraction)()
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening connection", err)
		return
	}

	for _, command := range cmds {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
		if err != nil {
			log.Fatal("Error creating slash command", err)
			return
		}
	}

	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	if err := discord.Close(); err != nil {
		log.Fatal("Error closing connection", err)
		return
	}
	log.Println("Bot is shutting down")
	saveSession()
}
