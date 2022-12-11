package handler

import (
	"github.com/DenChenn/yenno-dc/config"
	"github.com/DenChenn/yenno-dc/dao"
	"github.com/bwmarrin/discordgo"
	"github.com/teris-io/shortid"
)

type handler struct {
	DAO         *dao.DAO
	IDGenerator *shortid.Shortid
}

func NewHandler(dao *dao.DAO, generator *shortid.Shortid) *handler {
	return &handler{
		DAO:         dao,
		IDGenerator: generator,
	}
}

func (h *handler) Select(ic *discordgo.InteractionCreate) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if ic.Type == discordgo.InteractionApplicationCommand {
		switch ic.ApplicationCommandData().Name {
		case config.CreateDeploymentConfigCommand:
			return h.CreateDeploymentConfig
		case config.ListAllDeploymentConfigCommand:
			return h.ListAllDeploymentConfig
		case config.DeleteDeploymentConfigCommand:
			return h.DeleteDeploymentConfig
		case config.DeployWithDeploymentConfigCommand:
			return h.DeployWithDeploymentConfig
		case config.GetDeploymentConfigYamlCommand:
			return h.GetDeploymentConfigYaml
		}
	} else if ic.Type == discordgo.InteractionModalSubmit {
		switch ic.ModalSubmitData().CustomID {
		case config.CreateDeploymentConfigCommand:
			return h.ReceiveCreateDeploymentConfig
		case config.DeleteDeploymentConfigCommand:
			return h.ReceiveDeleteDeploymentConfig
		case config.DeployWithDeploymentConfigCommand:
			return h.ReceiveDeployWithDeploymentConfig
		case config.GetDeploymentConfigYamlCommand:
			return h.ReceiveGetDeploymentConfigYaml
		}
	}
	return nil
}
