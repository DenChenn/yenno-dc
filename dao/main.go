package dao

import "github.com/go-redis/redis/v9"

type DAO struct {
	DeploymentConfig *DeploymentConfigDAO
}

func New(client *redis.Client) *DAO {
	return &DAO{
		DeploymentConfig: NewDeploymentConfigDAO(client),
	}
}
