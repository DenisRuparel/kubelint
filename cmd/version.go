package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v3.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show KubeLint version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("KubeLint Version:", Version)
	},
}