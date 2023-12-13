package common

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mission-focused/mk8s/src/types"
	"gopkg.in/yaml.v3"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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

// SSHConfig represents the configuration for connecting to an SSH server
type SSHConfig struct {
	Host       string
	Port       int
	User       string
	Password   string // You can use a private key instead of a password for more security
	PrivateKey string
}

// RunCommand executes the given command on the remote machine and returns the output
func RunCommand(sshConfig SSHConfig, command string) (string, error) {
	// Build the SSH client configuration
	config := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConfig.Password),
			// You can also use ssh.PublicKeys(privateKey) for key-based authentication
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use a proper HostKeyCallback for production
	}

	// Connect to the SSH server
	addr := fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("Failed to dial: %v", err)
	}
	defer client.Close()

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create session: %v", err)
	}
	defer session.Close()

	// Capture the output of the command
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("Failed to run command: %v", err)
	}

	return string(output), nil
}

// CopyFileWithSSH copies a file from source to destination using SFTP over SSH
func CopyFileWithSSH(sshConfig SSHConfig, sourceFilePath, destFilePath string) error {
	// Read private key file
	privateKey, err := os.ReadFile(sshConfig.PrivateKey)
	if err != nil {
		return err
	}

	// Create a signer from the private key
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return err
	}

	// Connect to the SSH server
	config := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Initiate an SFTP session
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// Open source file on the local machine
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file on the remote server
	destFile, err := sftpClient.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy file contents to the destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	fmt.Println("File copied successfully!")
	return nil
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
