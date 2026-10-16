package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/DenisRuparel/kubelint/internal/builder"
	"github.com/DenisRuparel/kubelint/internal/scanner"
	"github.com/DenisRuparel/kubelint/internal/utils"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
)

var applyValues string

var applyCmd = &cobra.Command{
	Use:   "apply [project-path]",
	Short: "Render and apply Kubernetes manifests",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]

		fmt.Println("Applying KubeLint manifests...")

		output, err := builder.Build(projectPath, applyValues)
		if err != nil {
			return fmt.Errorf("build failed before apply: %w", err)
		}

		scanResult := scanner.ScanRenderedYAML(output)
		summary := scanResult.Summary

		fmt.Println("\n🔍 Apply Summary")
		fmt.Println("---------------------------------")
		fmt.Printf("CRITICAL : %d\n", summary.Critical)
		fmt.Printf("WARNING  : %d\n", summary.Warning)
		fmt.Printf("INFO     : %d\n", summary.Info)
		fmt.Println("---------------------------------")

		if summary.Critical > 0 {
			fmt.Println("\n❌ Critical Issues:")
			fmt.Println("---------------------------------")

			for _, issue := range scanResult.Issues {
				if issue.Severity == validators.Critical {
					fmt.Printf("• %s\n", issue.Message)
				}
			}

			fmt.Println("---------------------------------")
			fmt.Println("\n❌ Deployment blocked due to critical issues.")
			return fmt.Errorf("deployment blocked: %d critical issue(s) detected", summary.Critical)
		}

		// kubectl apply -f -
		kubectl := exec.Command("kubectl", "apply", "-f", "-")
		clean := utils.CleanYAML(output)
		kubectl.Stdin = strings.NewReader(clean)

		out, err := kubectl.CombinedOutput()
		if err != nil {
			fmt.Println("\n❌ kubectl apply failed:")
			fmt.Println("---------------------------------")
			fmt.Println(string(out))
			fmt.Println("---------------------------------")
			fmt.Println("💡 Tip: Run `kubelint build` to see which resource is broken")
			return fmt.Errorf("kubectl apply failed: %w", err)
		}

		fmt.Println("✔ Applied successfully")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyValues, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(applyCmd)
}
