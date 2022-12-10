package dao

import (
	"context"
	"encoding/json"
	"github.com/DenChenn/yenno-dc/model"
	"github.com/go-redis/redis/v9"
)

type DeploymentConfigDAO struct {
	RedisClient *redis.Client
}

func NewDeploymentConfigDAO(redisClient *redis.Client) *DeploymentConfigDAO {
	return &DeploymentConfigDAO{
		RedisClient: redisClient,
	}
}

func (d *DeploymentConfigDAO) Get(ctx context.Context, id string) (*model.DeploymentConfig, error) {
	rawDeploymentConfig, err := d.RedisClient.Get(ctx, id).Bytes()
	if err != nil {
		return nil, err
	}

	var deploymentConfig model.DeploymentConfig
	if err := json.Unmarshal(rawDeploymentConfig, &deploymentConfig); err != nil {
		return nil, err
	}
	return &deploymentConfig, nil
}

func (d *DeploymentConfigDAO) GetAll(ctx context.Context) ([]*model.DeploymentConfig, error) {
	var deploymentConfigs []*model.DeploymentConfig
	keys, err := d.RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		deploymentConfig, err := d.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		deploymentConfigs = append(deploymentConfigs, deploymentConfig)
	}
	return deploymentConfigs, nil
}

func (d *DeploymentConfigDAO) Create(ctx context.Context, deploymentConfig *model.DeploymentConfig) error {
	rawDeploymentConfig, err := json.Marshal(deploymentConfig)
	if err != nil {
		return err
	}

	if err := d.RedisClient.Set(ctx, deploymentConfig.ID, rawDeploymentConfig, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (d *DeploymentConfigDAO) Delete(ctx context.Context, id string) error {
	if err := d.RedisClient.Del(ctx, id).Err(); err != nil {
		return err
	}
	return nil
}
