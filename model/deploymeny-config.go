package model

import "time"

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeploymentConfig struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	LimitCPU      int32      `json:"limit_cpu"`
	RequestCPU    int32      `json:"request_cpu"`
	LimitMemory   int32      `json:"limit_memory"`
	RequestMemory int32      `json:"request_memory"`
	ImageURL      string     `json:"image_url"`
	Node          string     `json:"node"`
	ContainerPort int32      `json:"container_port"`
	Env           []Env      `json:"env"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
