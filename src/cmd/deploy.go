/*
Copyright © 2023 MissionFocusedDevelopers
*/
package cmd

import (
	"fmt"

	"github.com/brandtkeller/mk8s/src/internal/common"
	"github.com/spf13/cobra"
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
		config, err := common.ConfigFromPath(path)
		if err != nil {
			return err
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
