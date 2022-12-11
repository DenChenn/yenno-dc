package handler

import (
	"context"
	"log"

	"github.com/DenChenn/yenno-dc/config"
	"github.com/bwmarrin/discordgo"
)

func (h *handler) CreateDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: config.CreateDeploymentConfigCommand,
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

func (h *handler) GetDeploymentConfigYaml(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: config.GetDeploymentConfigYamlCommand,
			Title:    "Get deployment config yaml",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "Enter deployment config id to get yaml",
							Label:    "Enter deployment config id to get yaml",
							Style:    discordgo.TextInputShort,
							Required: true,
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
	deploymentConfigList, err := h.DAO.DeploymentConfig.GetAll(context.Background())
	if err != nil {
		log.Println(err)
	}
	var fd []*discordgo.MessageEmbedField
	for _, deploymentConfig := range deploymentConfigList {
		fd = append(fd, &discordgo.MessageEmbedField{
			Name:  deploymentConfig.Name,
			Value: deploymentConfig.ID,
		})
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  "📜 All deployment config",
					Color:  config.White,
					Fields: fd,
				},
			},
		},
	}); err != nil {
		log.Println(err)
	}
}

func (h *handler) DeleteDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: config.DeleteDeploymentConfigCommand,
			Title:    "Delete deployment config",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "Enter deployment config id to delete",
							Label:    "Enter deployment config id to delete",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
			},
		},
	}); err != nil {
		log.Println(err)
	}
}

func (h *handler) DeployWithDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: config.DeployWithDeploymentConfigCommand,
			Title:    "Deploy with deployment config",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "Enter deployment config id to deploy with",
							Label:    "Enter deployment config id to deploy with",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
			},
		},
	}); err != nil {
		log.Println(err)
	}
}
