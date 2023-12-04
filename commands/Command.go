package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	ExecuteSlash(session *discordgo.Session, interaction *discordgo.InteractionCreate)
	ExecuteDash(session *discordgo.Session, messageCreate *discordgo.MessageCreate, args string)
}
