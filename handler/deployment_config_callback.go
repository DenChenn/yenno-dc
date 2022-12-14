package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"

	"github.com/DenChenn/yenno-dc/config"
	"github.com/DenChenn/yenno-dc/model"
	"github.com/DenChenn/yenno-dc/template"
	"github.com/bwmarrin/discordgo"
)

func (h *handler) ReceiveCreateDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	deploymentConfig.RequestCPU = cpuMemorySlice[0]
	deploymentConfig.LimitCPU = cpuMemorySlice[1]
	deploymentConfig.RequestMemory = cpuMemorySlice[2]
	deploymentConfig.LimitMemory = cpuMemorySlice[3]

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
		Value: deploymentConfig.RequestCPU,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "LimitCPU",
		Value: deploymentConfig.LimitCPU,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "RequestMemory",
		Value: deploymentConfig.RequestMemory,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "LimitMemory",
		Value: deploymentConfig.LimitMemory,
	})
	fd = append(fd, &discordgo.MessageEmbedField{
		Name:  "Env",
		Value: data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	})

	_, err = s.ChannelMessageSendEmbed(config.Env.ChannelId, &discordgo.MessageEmbed{
		Title:  "??? Create Deployment Config Success",
		Color:  config.White,
		Fields: fd,
	})
	if err != nil {
		log.Println(err)
	}

	// generate k8s yaml
	yamlFilePath, err := template.GenerateYaml(deploymentConfig)
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(yamlFilePath)
	if err != nil {
		log.Println(err)
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "???? Here is your yaml file",
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
	template.RemoveYaml(yamlFilePath)
}

func (h *handler) ReceiveDeleteDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// get deployment config id
	data := i.ModalSubmitData()
	deploymentConfigID := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	// delete deployment config
	if err := h.DAO.DeploymentConfig.Delete(context.Background(), deploymentConfigID); err != nil {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Delete Deployment Config Fail",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}

	// return success status
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "??? Delete Deployment Config Success",
		},
	}); err != nil {
		log.Println(err)
	}
}

func (h *handler) ReceiveDeployWithDeploymentConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: replace hardcoded variables
	host := "140.113.179.35:2030"
	user := "master"
	pKey := []byte(config.Env.AisfMasterPrivateKey)
	serverDeploymentFileRoot := "/home/master"

	// get deployment config id
	data := i.ModalSubmitData()
	deploymentConfigID := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	// get deployment config
	deploymentConfig, err := h.DAO.DeploymentConfig.Get(context.Background(), deploymentConfigID)
	if err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to get deployment config from database",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(pKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// SSH client config
	clientConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	sshClient, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to connect to server over ssh",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}
	defer sshClient.Close()

	// Create a new SCP client
	scpClient, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to connect to server over ssh",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}

	// Connect to the remote server
	if err := scpClient.Connect(); err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to connect to server over ssh",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}
	defer scpClient.Close()

	// Upload yaml file
	yamlFilePath, err := template.GenerateYaml(deploymentConfig)
	if err != nil {
		log.Println(err)
	}
	file, err := os.Open(yamlFilePath)
	if err != nil {
		log.Println(err)
	}
	filename := deploymentConfig.Name + "_" + deploymentConfig.ID + ".yaml"
	remoteFilePath := filepath.Join(serverDeploymentFileRoot, filename)
	permission := "0655"
	if err := scpClient.CopyFromFile(context.Background(), *file, remoteFilePath, permission); err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to copy file to server",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}
	template.RemoveYaml(yamlFilePath)

	// create ssh session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to connect to server over ssh",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}
	defer session.Close()

	// kubectl apply file
	deployCmd := "kubectl apply -f " + remoteFilePath
	removeDeploymentFileCmd := "rm " + remoteFilePath
	if _, err := session.CombinedOutput(deployCmd + " && " + removeDeploymentFileCmd); err != nil {
		log.Println(err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Failed to apply deployment config",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "??? Successfully deployed config ID: " + deploymentConfigID,
		},
	}); err != nil {
		log.Println(err)
	}
	// os.Remove(resultFilePath)
}

func (h *handler) ReceiveGetDeploymentConfigYaml(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// get deployment config id
	data := i.ModalSubmitData()
	deploymentConfigID := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	// delete deployment config
	deploymentConfig, err := h.DAO.DeploymentConfig.Get(context.Background(), deploymentConfigID)
	if err != nil {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "??? Get Deployment Config Yaml Fail",
			},
		}); err != nil {
			log.Println(err)
		}
		return
	}

	// generate k8s yaml
	yamlFilePath, err := template.GenerateYaml(deploymentConfig)
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(yamlFilePath)
	if err != nil {
		log.Println(err)
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "???? Here is your yaml file",
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
	template.RemoveYaml(yamlFilePath)
}
