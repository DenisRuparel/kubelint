package cmd

import (
	"fmt"
	"github.com/DenisRuparel/kubelint/internal/builder"
	"github.com/DenisRuparel/kubelint/internal/scanner"
	"github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var applyValues string

var applyCmd = &cobra.Command{
	Use:   "apply [project-path]",
	Short: "Render and apply Kubernetes manifests",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		projectPath := args[0]

		fmt.Println("Applying KubeLint manifests...")

		output, err := builder.Build(projectPath, applyValues)

		if err != nil {
			fmt.Println("❌", err)
			os.Exit(1)
		}

		scanResult := scanner.ScanRenderedYAML(output)
		summary := scanResult.Summary

		// fmt.Println("DEBUG APPLY OUTPUT:\n", output)

		fmt.Println("\n🔍 Scan Summary")
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
			fmt.Println("❌ Deployment blocked due to critical issues.")
			os.Exit(1)
		}

		if err != nil {
			fmt.Println("❌", err)
			os.Exit(1)
		}

		// kubectl apply -f -
		kubectl := exec.Command("kubectl", "apply", "-f", "-")
		kubectl.Stdin = strings.NewReader(output)

		out, err := kubectl.CombinedOutput()
		if err != nil {
			fmt.Println("❌ kubectl apply failed:")
			fmt.Println(string(out))
			os.Exit(1)
		}

		fmt.Println("✔ Applied successfully")
		fmt.Println(string(out))
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyValues, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(applyCmd)
}
