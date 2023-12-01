/*
Copyright © 2023 MissionFocusedDevelopers
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download artifacts for kubernetes cluster installation",
	Long:  `download artifacts for kubernetes cluster installation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("download called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
