package common

import (
	"bytes"
	"errors"
	"fmt"
	// "io"
	"net"
	"os"
	"runtime"
	"strconv"

	operator "github.com/alexellis/k3sup/pkg/operator"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

type SSHConfig struct {
	Host       string
	Port       int
	User       string
	Password   string // You can use a private key instead of a password for more security
	PrivateKey string
}

// CopyFileWithSSH copies a file from source to destination using SFTP over SSH
func (sshConfig SSHConfig) CopyFileWithSSH(sourceFilePath, destFilePath string) error {
	// Read private key file
	privateKey, err := os.ReadFile(expandPath(sshConfig.PrivateKey))
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
		return fmt.Errorf("Error creating destination file: %s with error: %s", destFilePath, err)
	}
	defer destFile.Close()

	// Copy file contents to the destination
	_, err = destFile.ReadFrom(sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func (sshConfig SSHConfig) RunRemoteCommand(command string) (result operator.CommandRes, err error) {

	sshOperator, sshOperatorDone, errored, err := connectOperator(sshConfig.User, sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), sshConfig.PrivateKey)
	if errored {
		return result, err
	}

	if sshOperatorDone != nil {

		defer sshOperatorDone()
	}

	result, err = sshOperator.Execute(command)

	if err != nil {
		return result, fmt.Errorf("error received processing command: %s", err)
	}

	// fmt.Printf("Result: %s %s\n", string(res.StdOut), string(res.StdErr))

	// if err = obtainKubeconfig(sshOperator, getConfigcommand, host, context, localKubeconfig, merge, printConfig); err != nil {
	// 	return err
	// }

	return result, err

}

type DoneFunc func()

func connectOperator(user string, address string, sshKeyPath string) (*operator.SSHOperator, DoneFunc, bool, error) {
	var sshOperator *operator.SSHOperator
	var initialSSHErr error
	var closeSSHAgentFunc func() error

	doneFunc := func() {
		if sshOperator != nil {
			sshOperator.Close()
		}
		if closeSSHAgentFunc != nil {
			closeSSHAgentFunc()
		}
	}

	if runtime.GOOS != "windows" {
		var sshAgentAuthMethod ssh.AuthMethod
		sshAgentAuthMethod, initialSSHErr = sshAgentOnly()
		if initialSSHErr == nil {

			config := &ssh.ClientConfig{
				User:            user,
				Auth:            []ssh.AuthMethod{sshAgentAuthMethod},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			sshOperator, initialSSHErr = operator.NewSSHOperator(address, config)
		}
	} else {
		initialSSHErr = errors.New("ssh-agent unsupported on windows")
	}

	if initialSSHErr != nil {
		path := expandPath(sshKeyPath)
		publicKeyFileAuth, closeSSHAgent, err := loadPublickey(path)
		if err != nil {
			return nil, nil, true, fmt.Errorf("unable to load the ssh key with path %q: %w", path, err)
		}

		defer closeSSHAgent()

		config := &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{publicKeyFileAuth},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		sshOperator, err = operator.NewSSHOperator(address, config)
		if err != nil {
			return nil, nil, true, fmt.Errorf("unable to connect to %s over ssh: %w", address, err)
		}
	}

	return sshOperator, doneFunc, false, nil
}

func sshAgentOnly() (ssh.AuthMethod, error) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
}

func expandPath(path string) string {
	res, _ := homedir.Expand(path)
	return res
}

func sshAgent(publicKeyPath string) (ssh.AuthMethod, func() error) {
	if sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		sshAgent := agent.NewClient(sshAgentConn)

		keys, _ := sshAgent.List()
		if len(keys) == 0 {
			return nil, sshAgentConn.Close
		}

		pubkey, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, sshAgentConn.Close
		}

		authkey, _, _, _, err := ssh.ParseAuthorizedKey(pubkey)
		if err != nil {
			return nil, sshAgentConn.Close
		}
		parsedkey := authkey.Marshal()

		for _, key := range keys {
			if bytes.Equal(key.Blob, parsedkey) {
				return ssh.PublicKeysCallback(sshAgent.Signers), sshAgentConn.Close
			}
		}
	}
	return nil, func() error { return nil }
}

func loadPublickey(path string) (ssh.AuthMethod, func() error, error) {
	noopCloseFunc := func() error { return nil }

	key, err := os.ReadFile(path)
	if err != nil {
		return nil, noopCloseFunc, fmt.Errorf("unable to read file: %s, %s", path, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		if _, ok := err.(*ssh.PassphraseMissingError); !ok {
			return nil, noopCloseFunc, fmt.Errorf("unable to parse private key: %s", err.Error())
		}

		agent, close := sshAgent(path + ".pub")
		if agent != nil {
			return agent, close, nil
		}

		defer close()

		fmt.Printf("Enter passphrase for '%s': ", path)
		STDIN := int(os.Stdin.Fd())
		bytePassword, _ := term.ReadPassword(STDIN)

		// Ignore any error from reading stdin to retain existing behaviour for unit test in
		// install_test.go

		// if err != nil {
		// 	return nil, noopCloseFunc, fmt.Errorf("reading password from stdin failed: %s", err.Error())
		// }

		fmt.Println()

		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, bytePassword)
		if err != nil {
			return nil, noopCloseFunc, fmt.Errorf("parse private key with passphrase failed: %s", err)
		}
	}

	return ssh.PublicKeys(signer), noopCloseFunc, nil
}
