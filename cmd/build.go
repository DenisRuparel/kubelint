package cmd

import (
	"fmt"

	"github.com/DenisRuparel/kubelint/internal/builder"
	"github.com/DenisRuparel/kubelint/internal/scanner"
	"github.com/DenisRuparel/kubelint/internal/utils"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
)

var valuesFile string

var buildCmd = &cobra.Command{
	Use:   "build [project-path]",
	Short: "Render Kubernetes templates with values",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]

		output, err := builder.Build(projectPath, valuesFile)
		if err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		// Validate rendered YAML
		scanResult := scanner.ScanRenderedYAML(output)
		summary := scanResult.Summary

		if summary.Critical > 0 {
			fmt.Println("\n🔍 Build Validation")
			fmt.Println("---------------------------------")
			fmt.Printf("CRITICAL : %d\n", summary.Critical)
			fmt.Printf("WARNING  : %d\n", summary.Warning)
			fmt.Printf("INFO     : %d\n", summary.Info)
			fmt.Println("---------------------------------")
			fmt.Println("\n⚠️ Build Issues detected:")
			fmt.Println("---------------------------------")

			for _, issue := range scanResult.Issues {
				if issue.Severity == validators.Critical {
					fmt.Printf("• %s\n", issue.Message)
				}
			}

			fmt.Println("---------------------------------")
			fmt.Println("Status: FAILED ❌")
			return fmt.Errorf("build validation failed with %d critical issue(s)", summary.Critical)
		}

		clean := utils.CleanYAML(output)
		fmt.Print(clean)
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&valuesFile, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(buildCmd)
}
