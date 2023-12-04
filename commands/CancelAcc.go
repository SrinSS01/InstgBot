package commands

import "github.com/bwmarrin/discordgo"

type CancelAccCommand struct {
	Command *discordgo.ApplicationCommand
}

var CancelAcc = AccCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "cancelacc",
		Description: "cancel the process of account creation",
	},
}

func (ca *CancelAccCommand) ExecuteDash(s *discordgo.Session, m *discordgo.MessageCreate, _ string) {
	authorId := m.Author.ID
	if authorId == s.State.User.ID {
		return
	}
	key := authorId + m.GuildID
	i := InfoMap[key]
	if i == nil {
		return
	}
	_, _ = s.ChannelMessageSendReply(m.ChannelID, "Cancelling process...", m.Reference())
	delete(InfoMap, key)
}

func (ca *CancelAccCommand) ExecuteSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	m := i.Message
	authorId := m.Author.ID
	if authorId == s.State.User.ID {
		return
	}
	key := authorId + m.GuildID
	inf := InfoMap[key]
	if inf == nil {
		return
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Cancelling process...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	delete(InfoMap, key)
}
