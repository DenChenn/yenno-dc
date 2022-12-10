package handler

import (
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
		case "create-deployment-config":
			return h.CreateDeploymentConfig
		case "list-all-deployment-config":
			return h.ListAllDeploymentConfig
		case "delete-deployment-config":
			return h.DeleteDeploymentConfig
		case "deploy-with-config":
			return h.DeployWithConfig
		}
	} else if ic.Type == discordgo.InteractionModalSubmit {
		switch ic.ModalSubmitData().CustomID {
		case "create-deployment-config":
			return h.ReceiveCreateConfig
		}
	}
	return nil
}
