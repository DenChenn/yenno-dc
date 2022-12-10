package handler

import (
	"context"
	"fmt"
	"github.com/DenChenn/yenno-dc/config"
	"github.com/DenChenn/yenno-dc/model"
	"github.com/DenChenn/yenno-dc/template"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (h *handler) ReceiveCreateConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	uuid, err := h.IDGenerator.Generate()
	if err != nil {
		log.Println(err)
	}
	currentTime := time.Now()
	deploymentConfig := &model.DeploymentConfig{
		ID:        uuid,
		Name:      data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		ImageURL:  data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	// precess format of <request_X>|<limit_X>|<request_Y>|<limit_Y>
	cpuMemorySlice := strings.Split(
		data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		"|",
	)
	requestCPU, _ := strconv.Atoi(cpuMemorySlice[0])
	limitCPU, _ := strconv.Atoi(cpuMemorySlice[1])
	requestMemory, _ := strconv.Atoi(cpuMemorySlice[2])
	limitMemory, _ := strconv.Atoi(cpuMemorySlice[3])

	deploymentConfig.RequestCPU = int32(requestCPU)
	deploymentConfig.LimitCPU = int32(limitCPU)
	deploymentConfig.RequestMemory = int32(requestMemory)
	deploymentConfig.LimitMemory = int32(limitMemory)

	// process format of <node>|<container_port>
	nodeContainerPortSlice := strings.Split(
		data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		"|",
	)
	containerPort, _ := strconv.Atoi(nodeContainerPortSlice[1])
	deploymentConfig.Node = nodeContainerPortSlice[0]
	deploymentConfig.ContainerPort = int32(containerPort)

	// precess env variables
	linesOfEnv := strings.Split(
		data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
		"\n",
	)
	deploymentConfig.Env = make([]model.Env, 0)
	for _, line := range linesOfEnv {
		if line == "" {
			continue
		}
		env := strings.Split(line, "=")
		deploymentConfig.Env = append(deploymentConfig.Env, model.Env{
			Key:   env[0],
			Value: line[len(env[0])+1:],
		})
	}

	if err := h.DAO.DeploymentConfig.Create(context.Background(), deploymentConfig); err != nil {
		log.Println(err)
	}

	// return success status
	var fd []*discordgo.MessageEmbedField
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "ID",
		Value: deploymentConfig.ID,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "Name",
		Value: deploymentConfig.Name,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "ImageURL",
		Value: deploymentConfig.ImageURL,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "Node",
		Value: deploymentConfig.Node,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "ContainerPort",
		Value: strconv.Itoa(int(deploymentConfig.ContainerPort)),
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "RequestCPU",
		Value: fmt.Sprintf("%d", deploymentConfig.RequestCPU),
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "LimitCPU",
		Value: fmt.Sprintf("%d", deploymentConfig.LimitCPU),
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "RequestMemory",
		Value: fmt.Sprintf("%d", deploymentConfig.RequestMemory),
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "LimitMemory",
		Value: fmt.Sprintf("%d", deploymentConfig.LimitMemory),
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "Env",
		Value: data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	})

	_, err = s.ChannelMessageSendEmbed(config.Env.ChannelId, &discordgo.MessageEmbed{
		Title:  "Create Deployment Config Success ✅",
		Color:  16775936,
		Fields: fd,
	})
	if err != nil {
		log.Println(err)
	}

	// generate k8s manifest
	if err := template.GenerateManifest(deploymentConfig); err != nil {
		log.Println(err)
	}

	filename := deploymentConfig.Name + "_" + deploymentConfig.ID + ".yaml"
	manifestFilePath := filepath.Join(config.RootPath, filename)
	file, err := os.Open(manifestFilePath)
	if err != nil {
		log.Println(err)
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Here is your manifest file",
			Files: []*discordgo.File{
				{
					ContentType: "text/yaml",
					Name:        deploymentConfig.Name + "_" + deploymentConfig.ID + ".yaml",
					Reader:      file,
				},
			},
		},
	}); err != nil {
		log.Println(err)
	}
	template.RemoveManifest(manifestFilePath)
}
