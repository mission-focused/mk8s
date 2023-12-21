package distro

import (
	"fmt"
	"os"
	"strings"

	"github.com/mission-focused/mk8s/src/internal/common"
	"github.com/mission-focused/mk8s/src/types"
)

var (
	createConfigDir  = "sudo mkdir -p /etc/rancher/rke2"
	artifactDir      = "~/rke2-artifacts"
	configDest       = "/etc/rancher/rke2/config.yaml"
	installExecution = ""
)

// TODO: maybe separate this to another package in the future

// TODO: Here we can have more fun with concurrency in the future
func installMultiRKE2(config types.MultiConfig) (err error) {

	// Likely we want to establish some concurrency here in the future
	nodes, err := common.NodeMapFromSlice(config.Nodes)
	if err != nil {
		return err
	}

	artifacts, err := ValidateArtifacts(config)
	if err != nil {
		return err
	}
	// Prioritize installation on the primary node if identified
	// TODO: future - then run the install on all other nodes simultaneously
	if _, ok := nodes["primary"]; ok {
		// There should always only be one primary
		err = installRKE2(nodes["primary"][0], artifacts)
		if err != nil {
			return err
		}

		// grab and alter the kubeconfig here for use immediately
	}

	return nil
}

// Single node installation
func installRKE2(node types.NodeConfig, artifacts map[string]types.Artifact) error {
	//		Run the installation script
	//		INSTALL_RKE2_ARTIFACT_PATH=/home/dev/rke2-artifacts INSTALL_RKE2_TYPE='agent' sh /home/dev/rke2-artifacts/install.sh"

	//		Enable the rke2 service on the node
	//		Start the rke2 service on the node

	// Create the config files - retaining in the artifacts directory for reference
	ipString := strings.Replace(node.Address, ".", "-", -1)
	fileName := "artifacts/" + ipString + "-config.yaml"

	err := common.CreateConfigFile(node.Config, fileName)
	if err != nil {
		return err
	}

	if !node.Local {
		// remote installation - ssh required
		// Create the sshconfig
		sshconfig := common.SSHConfig{
			Host:       node.Address,
			Port:       22,
			User:       node.User,
			PrivateKey: node.SshKeyPath,
		}

		// Not expecting any result here
		_, err := sshconfig.RunRemoteCommand(createConfigDir)
		if err != nil {
			return err
		}

		createArtifactDir := fmt.Sprintf("sudo mkdir -p %s", artifactDir)
		_, err = sshconfig.RunRemoteCommand(createArtifactDir)
		if err != nil {
			return err
		}

		fmt.Println("Attempting to copy config file")
		// Copy config file to rke2 directory
		err = sshconfig.CopyFileWithSSH(fileName, "rke2-artifacts/"+ipString+"-config.yaml")
		if err != nil {
			return err
		}

		fmt.Println("Attempting to move config file to final destination")
		configMove := fmt.Sprintf("sudo cp %s %s", artifactDir+"/"+ipString+"-config.yaml", configDest)
		_, err = sshconfig.RunRemoteCommand(configMove)
		if err != nil {
			return err
		}

	} else {
		// Local installation - no ssh required

		// Create config directory
		err := common.ExecuteLocalCommand(createConfigDir)
		if err != nil {
			return err
		}

		configCopyCmd := fmt.Sprintf("sudo cp %s %s", fileName, configDest)
		// Copy config file to rke2 directory
		err = common.ExecuteLocalCommand(configCopyCmd)
		if err != nil {
			return err
		}

		// If we use local - we can use the artifacts directory as the install source
		// get PWD and execute the install script
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		installCmd := fmt.Sprintf("sudo INSTALL_RKE2_ARTIFACT_PATH=%s INSTALL_RKE2_TYPE='%s' sh %s", pwd+"/artifacts", node.Role, pwd+"/artifacts/install.sh")

		err = common.ExecuteLocalCommand(installCmd)
		if err != nil {
			return err
		}

		systemctlCmd := fmt.Sprintf("sudo systemctl enable rke2-%s.service && sudo systemctl start rke2-%s.service", node.Role, node.Role)
		err = common.ExecuteLocalCommand(systemctlCmd)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
