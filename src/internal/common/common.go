package common

import (
	"fmt"
	"os"

	"github.com/brandtkeller/mk8s/src/types"
	"gopkg.in/yaml.v3"
)

func ConfigFromPath(path string) (types.MultiConfig, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return types.MultiConfig{}, fmt.Errorf("Path: %v does not exist - unable to digest document\n", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return types.MultiConfig{}, err
	}
	// marshall to an object? object name?
	var config types.MultiConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return types.MultiConfig{}, fmt.Errorf("Error marshalling yaml: %s\n", err.Error())
	}

	return config, nil
}
