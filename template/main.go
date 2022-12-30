package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/DenChenn/yenno-dc/config"
	"github.com/DenChenn/yenno-dc/model"
)

func ApplyTemplate(tmplPath string, config map[string]interface{}, outputPath string) error {
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if err := t.Execute(f, config); err != nil {
		return err
	}

	return nil
}

func GenerateYaml(deploymentConfig *model.DeploymentConfig) (string, error) {
	// create config map according to module detail
	var envMap = map[string]string{}
	for _, env := range deploymentConfig.Env {
		envMap[env.Key] = env.Value
	}

	c := map[string]interface{}{
		"Name":          deploymentConfig.Name,
		"ImageURL":      deploymentConfig.ImageURL,
		"RequestCPU":    deploymentConfig.RequestCPU,
		"LimitCPU":      deploymentConfig.LimitCPU,
		"RequestMemory": deploymentConfig.RequestMemory,
		"LimitMemory":   deploymentConfig.LimitMemory,
		"Node":          deploymentConfig.Node,
		"ContainerPort": deploymentConfig.ContainerPort,
		"Env":           envMap,
	}

	templatePath := filepath.Join(config.RootPath, "template", "deploy.tmpl")
	filename := deploymentConfig.Name + "_" + deploymentConfig.ID + ".yaml"
	outputPath := filepath.Join(config.RootPath, filename)

	return outputPath, ApplyTemplate(templatePath, c, outputPath)
}

func RemoveYaml(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}
