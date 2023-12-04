package commands

import "github.com/bwmarrin/discordgo"

type AddSessionIdCommand struct {
	Command *discordgo.ApplicationCommand
}

var AddSessionId = AddSessionIdCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "add-session-id",
		Description: "Add session ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "session",
				Description: "Session ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

func (ac *AddSessionIdCommand) ExecuteSlash(s *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ApplicationCommandData()
	sessionId := data.Options[0].StringValue()
	_ = s.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Saved session ID",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	SessionObj.Sessions = append(SessionObj.Sessions, sessionId)
}
