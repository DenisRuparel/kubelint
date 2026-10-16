package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubelint",
	Short: "KubeLint is a Kubernetes manifest linting CLI",
	Long: `KubeLint is a production-grade Kubernetes linting tool
that validates YAML manifests using best practices,
security checks, and deployment standards.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
