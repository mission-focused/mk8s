/*
Copyright © 2023 MissionFocusedDevelopers
*/
package cmd

import (
	"fmt"

	"github.com/brandtkeller/mk8s/src/internal/common"
	"github.com/brandtkeller/mk8s/src/internal/distro"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download artifacts for kubernetes cluster installation",
	Long:  `download artifacts for kubernetes cluster installation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("download called")

		path := args[0]
		// Read the manifest file - should be the first argument
		config, err := common.ConfigFromPath(path)
		if err != nil {
			return err
		}

		fmt.Println(config.Distro)

		// Download the artifacts
		err = distro.DownloadArtifacts(config)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
