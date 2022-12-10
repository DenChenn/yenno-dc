package handler

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func (h *handler) CreateDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "create-deployment-config",
			Title:    "Enter deployment config",
			Components: []discordgo.MessageComponent{
				// max length of row is five
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "name",
							Label:    "name",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "RequestCPU & LimitCPU & RequestMemory & LimitMemory",
							Label:       "CPU & Memory",
							Placeholder: "Enter in this format: <request_cpu>|<limit_cpu>|<request_memory>|<limit_memory>",
							Style:       discordgo.TextInputParagraph,
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "Image url",
							Label:    "Image url",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "Node & Container Port",
							Label:       "Node & Container Port",
							Placeholder: "Enter in this format: <node>|<container_port>",
							Style:       discordgo.TextInputParagraph,
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "Env",
							Label:       "Env",
							Placeholder: "hint: copy and paste .env file",
							Style:       discordgo.TextInputParagraph,
							Required:    true,
						},
					},
				},
			},
		},
	}); err != nil {
		log.Println(err)
	}
}
func (h *handler) ListAllDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}); err != nil {
		log.Println(err)
	}
}
func (h *handler) DeleteDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}); err != nil {
		log.Println(err)
	}
}
func (h *handler) DeployWithConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	}); err != nil {
		log.Println(err)
	}
}
