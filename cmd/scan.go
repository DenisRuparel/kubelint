package cmd

import (
	"fmt"
	"os"

	"github.com/DenisRuparel/kubelint/internal/scanner"
	validators "github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [file|files|directory]",
	Short: "Scan Kubernetes YAML files for issues",
	Args:  cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting KubeLint scan...")

		yamlFiles, err := scanner.CollectYAMLFiles(args)
		if err != nil {
			return fmt.Errorf("failed to collect YAML files: %w", err)
		}

		criticalCount := 0
		warningCount := 0
		infoCount := 0

		for _, file := range yamlFiles {
			// Read file ONCE from disk (fixes bug #4: 7-9 repeated file reads)
			data, err := os.ReadFile(file)
			if err != nil {
				criticalCount++
				fmt.Printf("\nScanning: %s\n", file)
				fmt.Printf("[CRITICAL] Failed to read file: %v\n", err)
				continue
			}

			var results []validators.ValidationResult

			// YAML syntax validation
			syntaxResult := validators.ValidateYAMLSyntaxBytes(data)
			syntaxResult.File = file
			results = append(results, syntaxResult)

			// Continue with deep validation only if syntax is valid
			if syntaxResult.Severity != validators.Critical {
				appendResults := func(res []validators.ValidationResult) {
					for _, r := range res {
						r.File = file
						results = append(results, r)
					}
				}

				appendResults(validators.ValidateCommonBytes(data))

				structRes := validators.ValidateStructureBytes(data)
				structRes.File = file
				results = append(results, structRes)

				appendResults(validators.ValidateDeploymentBytes(data))
				appendResults(validators.ValidateSecurityBytes(data))
				appendResults(validators.ValidateServiceBytes(data))
				appendResults(validators.ValidateConfigMapBytes(data))
				appendResults(validators.ValidateSecretBytes(data))
				appendResults(validators.ValidateIngressBytes(data))
				appendResults(validators.ValidateNamespacePolicyBytes(data))
			}

			fmt.Println()
			fmt.Printf("Scanning: %s\n", file)
			fmt.Println()

			fileScore := 100

			for _, result := range results {
				switch result.Severity {
				case validators.Critical:
					fmt.Printf("[CRITICAL] %s\n", result.Message)
					criticalCount++
					fileScore -= 20

				case validators.Warning:
					fmt.Printf("[WARNING]  %s\n", result.Message)
					warningCount++
					fileScore -= 10

				case validators.Info:
					fmt.Printf("[INFO]     %s\n", result.Message)
					infoCount++
				}
			}

			if fileScore < 0 {
				fileScore = 0
			}
		}

		overallScore := 100 - (criticalCount * 20) - (warningCount * 10)
		if overallScore < 0 {
			overallScore = 0
		}

		status := "EXCELLENT"
		if overallScore < 90 {
			status = "GOOD"
		}
		if overallScore < 70 {
			status = "NEEDS IMPROVEMENT"
		}
		if overallScore < 50 {
			status = "CRITICAL"
		}

		if criticalCount == 0 {
			fmt.Println()
			fmt.Println("------------ FINAL SUMMARY ------------")
			fmt.Printf("CRITICAL : %d\n", criticalCount)
			fmt.Printf("WARNING  : %d\n", warningCount)
			fmt.Printf("INFO     : %d\n", infoCount)
			fmt.Printf("\nOverall Score : %d/100\n", overallScore)
			fmt.Printf("Status        : %s\n", status)
			fmt.Println("---------------------------------------")
			return nil
		}

		fmt.Println("\n------------ SCAN FAILED ------------")
		fmt.Printf("CRITICAL : %d\n", criticalCount)
		fmt.Println("\nPlease resolve critical issues before proceeding.")
		fmt.Println("---------------------------------------")
		return fmt.Errorf("scan failed with %d critical issue(s)", criticalCount)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
