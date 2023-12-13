/*
Copyright © 2023 MissionFocusedDevelopers
*/
package cmd

import (
	"fmt"

	"github.com/mission-focused/mk8s/src/internal/common"
	"github.com/mission-focused/mk8s/src/internal/distro"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a kubernetes cluster",
	Long:  `install a kubernetes cluster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("install called")

		path := args[0]
		// Read the manifest file - should be the first argument
		config, err := common.ConfigFromPath(path)
		if err != nil {
			return err
		}
		// Perform validation of the configuration
		// This is likely a per-distribution function

		// Install
		err = distro.Install(config)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
