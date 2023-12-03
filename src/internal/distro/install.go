package distro

import (
	"fmt"

	"github.com/brandtkeller/mk8s/src/types"
)

func Install(config types.MultiConfig) error {

	nodeMap, err := nodeMapFromSlice(config.Nodes)
	if err != nil {
		return err
	}

	// fmt.Println(nodeMap)
	err = installRKE2(nodeMap)
	if err != nil {
		return err
	}

	// NodeMap is now established, we can now begin installation based on the distro

	return nil
}

// For future concurrency purposes and general mapping
// Create a map from the slice of NodeConfigs to a map of NodeConfigs
// What we care about is Primary server -> Secondary servers -> Agents
func nodeMapFromSlice(nodes []types.NodeConfig) (map[string][]types.NodeConfig, error) {

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
