package cmd

import (
	"fmt"
	"github.com/DenisRuparel/kubelint/internal/utils"
	"github.com/DenisRuparel/kubelint/internal/loader"
	"github.com/DenisRuparel/kubelint/internal/renderer"
	"github.com/DenisRuparel/kubelint/internal/scanner"
	"github.com/DenisRuparel/kubelint/internal/validator"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var valuesFile string
var finalOutput string

var buildCmd = &cobra.Command{
	Use:   "build [project-path]",
	Short: "Render Kubernetes templates with values",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		finalOutput = ""

		projectPath := args[0]
		templateDir := filepath.Join(projectPath, "templates")

		// Resolve values file safely
		localValuesFile := valuesFile
		if localValuesFile == "" {
			localValuesFile = filepath.Join(templateDir, "values.yaml")
		}

		fmt.Println("Starting KubeLint build...")

		// Path checks
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			fmt.Println("❌ templates folder not found")
			return
		}

		if _, err := os.Stat(localValuesFile); os.IsNotExist(err) {
			fmt.Println("❌ values.yaml not found")
			return
		}

		values, err := loader.LoadValues(localValuesFile)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		files, err := os.ReadDir(templateDir)
		if err != nil {
			fmt.Println("❌ Failed to read templates directory")
			return
		}

		var errors []string
		renderedCount := 0

		// 🔥 DO NOT PRINT YAML HERE
		for _, file := range files {
			if file.IsDir() || file.Name() == "values.yaml" {
				continue
			}

			filePath := filepath.Join(templateDir, file.Name())

			rendered, err := renderer.RenderTemplate(filePath, values)
			if err != nil {
				errors = append(errors,
					fmt.Sprintf("[%s] %v", file.Name(), err),
				)
				continue
			}

			// attach file metadata
			rendered = fmt.Sprintf("# FILE: %s\n%s", filePath, rendered)

			renderedCount++

			// accumulate ONLY
			finalOutput += "---\n"
			finalOutput += rendered

			if len(rendered) > 0 && rendered[len(rendered)-1] != '\n' {
				finalOutput += "\n"
			}
		}

		// 🔍 VALIDATE
		scanResult := scanner.ScanRenderedYAML(finalOutput)
		summary := scanResult.Summary

		
		// ❌ FAILURE CASE
		if summary.Critical > 0 {
			fmt.Println("\n🔍 Build Validation")
			fmt.Println("---------------------------------")
			fmt.Printf("CRITICAL : %d\n", summary.Critical)
			fmt.Printf("WARNING  : %d\n", summary.Warning)
			fmt.Printf("INFO     : %d\n", summary.Info)
			fmt.Println("---------------------------------")
			fmt.Println()
			fmt.Println("\n⚠️ Build Issues detected:")
			fmt.Println("---------------------------------")

			for _, issue := range scanResult.Issues {
				if issue.Severity == validators.Critical {
					fmt.Printf("• %s\n", issue.Message)
				}
			}

			fmt.Println("---------------------------------")

			fmt.Println("\n🚀 Build Summary")
			fmt.Println("---------------------------------")
			fmt.Printf("Templates rendered : %d\n", renderedCount)
			fmt.Printf("Critical Issues    : %d\n", summary.Critical)
			fmt.Printf("Warnings           : %d\n", summary.Warning)
			fmt.Printf("Info               : %d\n", summary.Info)
			fmt.Println("Status             : FAILED ❌")
			fmt.Println("---------------------------------")

			os.Exit(1)
		}

		// ✅ SUCCESS CASE → NOW PRINT YAML
		// fmt.Println("\nRendering templates...")
		fmt.Println()
		clean := utils.CleanYAML(finalOutput)
		fmt.Print(clean)

		// fmt.Println("\n🚀 Build Summary")
		// fmt.Println("---------------------------------")
		// fmt.Printf("Templates rendered : %d\n", renderedCount)
		// fmt.Printf("Critical Issues    : %d\n", summary.Critical)
		// fmt.Printf("Warnings           : %d\n", summary.Warning)
		// fmt.Printf("Info               : %d\n", summary.Info)
		// fmt.Println("Status             : SUCCESS ✅")
		// fmt.Println("---------------------------------")
	},
}

func init() {
	buildCmd.Flags().StringVarP(&valuesFile, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(buildCmd)
}
