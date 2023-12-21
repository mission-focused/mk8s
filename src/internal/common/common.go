package common

import (
	"fmt"
	"os"

	"github.com/mission-focused/mk8s/src/types"
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

// For future concurrency purposes and general mapping
// Create a map from the slice of NodeConfigs to a map of NodeConfigs
// What we care about is Primary server -> Secondary servers -> Agents
func NodeMapFromSlice(nodes []types.NodeConfig) (map[string][]types.NodeConfig, error) {

	nodeMap := make(map[string][]types.NodeConfig)

	for _, node := range nodes {
		// nodeMap[node.Role] = append(nodeMap[node.Role], node)

		if node.Primary == true {
			//Check if there is already a primary identified
			if len(nodeMap["primary"]) > 0 {
				return nodeMap, fmt.Errorf("More than one Primary already identified")
			}
			nodeMap["primary"] = append(nodeMap["primary"], node)
		} else if node.Role == "server" {
			// this means this is not the primary server node
			nodeMap["server"] = append(nodeMap["server"], node)
		} else {
			// Otherwise an agent node
			nodeMap["agent"] = append(nodeMap["agent"], node)
		}
	}

	return nodeMap, nil
}
