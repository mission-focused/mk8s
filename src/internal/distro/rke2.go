package distro

import (
	"fmt"
	"os"
	"path/filepath"
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

		// TODO: implement some method to check if the file already exists on the remote AND the hash matches
		// TODO: add a progress bar for the transfer process
		for key, artifact := range artifacts {
			if key == "images" {
				fmt.Println("Transferring images - this make take some time....")
			} else {
				fmt.Printf("Transferring artifact: %s\n", key)
			}

			err = sshconfig.CopyFileWithSSH("artifacts/"+artifact.Name, "rke2-artifacts/"+artifact.Name)
			if err != nil {
				return err
			}
		}

		fmt.Println("Running install script")

		installCmd := fmt.Sprintf("sudo INSTALL_RKE2_ARTIFACT_PATH=%s INSTALL_RKE2_TYPE='%s' sh %s", "~/rke2-artifacts", node.Role, "~/rke2-artifacts/install.sh")
		_, err = sshconfig.RunRemoteCommand(installCmd)
		if err != nil {
			return err
		}

		// Keeping this command separate with future intent on monitoring a channel for concurrency
		fmt.Println("Enabling and Starting rke2 service")

		enableStartCmd := fmt.Sprintf("sudo systemctl enable rke2-%s.service && sudo systemctl start rke2-%s.service", node.Role, node.Role)
		_, err = sshconfig.RunRemoteCommand(enableStartCmd)
		if err != nil {
			return err
		}

		if node.Primary {
			fmt.Println("Creating local copy of kubeconfig file")
			// move kubeconfig file to location where we can copy it out without escalated privileges
			res, err := sshconfig.RunRemoteCommand("sudo cat /etc/rancher/rke2/rke2.yaml")
			if err != nil {
				return err
			}

			path, err := os.Getwd()
			if err != nil {
				return err
			}
			absPath, _ := filepath.Abs(path)

			kubeconfig := string(res.StdOut)
			fmt.Printf("\nKubeconfig file: \n")
			fmt.Println(kubeconfig)

			config := strings.ReplaceAll(kubeconfig, "127.0.0.1", node.Address)

			if err := os.WriteFile(absPath+"/rke2.yaml", []byte(config), 0600); err != nil {
				return err
			}

		}

		return nil
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

}
