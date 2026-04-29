package cmd

import (
	"fmt"

	"github.com/DenisRuparel/kubelint/internal/scanner"
	"github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [file|files|directory]",
	Short: "Scan Kubernetes YAML files for issues",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting KubeLint scan...")

		yamlFiles, err := scanner.CollectYAMLFiles(args)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		criticalCount := 0
		warningCount := 0
		infoCount := 0

		for _, file := range yamlFiles {
			var results []validators.ValidationResult

			// YAML syntax validation
			syntaxResult := validators.ValidateYAMLSyntax(file)
			results = append(results, syntaxResult)

			// Continue only if syntax is valid
			if syntaxResult.Severity != validators.Critical {
				// Step 2: Common validation for all resources
				commonResults := validators.ValidateCommon(file)
				results = append(results, commonResults...)

				// Step 3: Structure validation
				structureResult := validators.ValidateStructure(file)
				results = append(results, structureResult)

				// Step 4: Deployment-specific validation
				deploymentResults := validators.ValidateDeployment(file)
				results = append(results, deploymentResults...)

				// Step 5: Service-specific validation
				serviceResults := validators.ValidateService(file)
				results = append(results, serviceResults...)
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

			// fmt.Printf("\nScore: %d/100\n\n", fileScore)
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

		} else {
			fmt.Println("\n------------ SCAN FAILED ------------")
			fmt.Printf("CRITICAL : %d\n", criticalCount)

			fmt.Println("\nPlease resolve critical issues before proceeding.")
			fmt.Println("---------------------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
