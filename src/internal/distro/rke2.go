package distro

import (
	"fmt"
	"strings"

	"github.com/brandtkeller/mk8s/src/internal/common"
	"github.com/brandtkeller/mk8s/src/types"
)

// TODO: maybe separate this to another package in the future

// TODO: Here we can have more fun with concurrency in the future
func installMultiRKE2(nodes map[string][]types.NodeConfig) (err error) {

	// Likely we want to establish some concurrency here in the future

	// Prioritize installation on the primary node if identified
	// TODO: future - then run the install on all other nodes simultaneously
	if _, ok := nodes["primary"]; ok {
		// There should always only be one primary
		err = installRKE2(nodes["primary"][0])
		if err != nil {
			return err
		}

		// grab and alter the kubeconfig here for use immediately
	}

	return nil
}

// Single node installation
func installRKE2(node types.NodeConfig) error {

	//

	//		Run the installation script
	//		INSTALL_RKE2_ARTIFACT_PATH=/home/dev/rke2-artifacts INSTALL_RKE2_TYPE='agent' sh /home/dev/rke2-artifacts/install.sh"

	//		Enable the rke2 service on the node
	//		Start the rke2 service on the node

	if !node.Local {
		// remote installation - ssh required
		// 		Create local copy of the config file
		ipString := strings.Replace(node.Address, ".", "-", -1)
		fileName := "artifacts/" + ipString + "-config.yaml"

		err := common.CreateConfigFile(node.Config, fileName)
		if err != nil {
			return err
		}

		// Create the sshconfig
		sshconfig := common.SSHConfig{
			Host:       node.Address,
			Port:       22,
			User:       node.User,
			PrivateKey: node.SshKeyPath,
		}

		// create /etc/rancher/rke2/ directory on the node
		output, err := common.RunCommand(sshconfig, "sudo mkdir -p /etc/rancher/rke2")
		if err != nil {
			return err
		}
		fmt.Println(output)

		// create the config file at /etc/rancher/rke2/config.yaml using the node config
		err = common.CopyFileWithSSH(sshconfig, fileName, "/etc/rancher/rke2/config.yaml")
		if err != nil {
			return err
		}

		// TODO: check for existence of files on the node
		// For now (given local testing) - we'll copy them
		// TODO: do we need to create a new ssh session for each command? Or can we reuse the same session?

		// TODO: Start Here - copy all artifacts from artifacts directory besides config files.

	} else {
		// Local installation - no ssh required
	}

	return nil
}
