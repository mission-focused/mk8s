/*
Copyright © 2023 MissionFocusedDevelopers
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/brandtkeller/mk8s/src/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a kubernetes cluster",
	Long:  `deploy a kubernetes cluster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("deploy called")

		path := args[0]
		// Read the manifest file - should be the first argument
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("Path: %v does not exist - unable to digest document\n", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// marshall to an object? object name?
		var config types.MultiConfig
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return fmt.Errorf("Error marshalling yaml: %s\n", err.Error())
		}

		fmt.Println(config)

		// Perform validation of the configuration
		// This is likely a per-distribution function

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
