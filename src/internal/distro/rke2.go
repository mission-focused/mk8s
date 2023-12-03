package distro

import "github.com/brandtkeller/mk8s/src/types"

// TODO: maybe separate this to another package in the future

// TODO: Here we can have more fun with concurrency in the future
func installRKE2(nodes map[string][]types.NodeConfig) (err error) {

	// Likely we want to establish some concurrency here in the future

	// Prioritize installation on the primary node if identified
	// TODO: future - then run the install on all other nodes simultaneously
	// For all nodes:

	//		copy the required artifacts to all nodes - check for existence first

	//		create /etc/rancher/rke2/ directory on the node

	// 		create the config file at /etc/rancher/rke2/config.yaml using the node config

	//		Run the installation script
	//		INSTALL_RKE2_ARTIFACT_PATH=/home/dev/rke2-artifacts INSTALL_RKE2_TYPE='agent' sh /home/dev/rke2-artifacts/install.sh"

	//		Enable the rke2 service on the node
	//		Start the rke2 service on the node

	return nil
}
