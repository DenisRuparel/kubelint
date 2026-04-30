package cmd

import (
	"fmt"
	"github.com/DenisRuparel/kubelint/internal/loader"
	"github.com/DenisRuparel/kubelint/internal/renderer"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var valuesFile string

var buildCmd = &cobra.Command{
	Use:   "build [project-path]",
	Short: "Render Kubernetes templates with values",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Println("\nRendering templates...")

		fmt.Println()

		var errors []string
		renderedCount := 0

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

			renderedCount++

			fmt.Println("---")
			fmt.Print(rendered)

			if len(rendered) > 0 && rendered[len(rendered)-1] != '\n' {
				fmt.Println()
			}
		}

		// 🔥 Print build issues (NEW)
		if len(errors) > 0 {
			fmt.Println("\n⚠️ Build Issues:")
			fmt.Println("---------------------------------")
			for _, e := range errors {
				fmt.Printf("• %s\n", e)
			}
			fmt.Println("---------------------------------")
			fmt.Printf("Total Errors: %d\n", len(errors))
		}

		// 🔥 Final summary (FIXED)
		fmt.Println("\n🚀 Build Summary")
		fmt.Println("---------------------------------")
		fmt.Printf("Templates rendered : %d\n", renderedCount)
		fmt.Printf("Errors             : %d\n", len(errors))

		if len(errors) == 0 {
			fmt.Println("Status             : SUCCESS ✅")
		} else {
			fmt.Println("Status             : PARTIAL ⚠️")
		}

		fmt.Println("---------------------------------")
	},
}

func init() {
	buildCmd.Flags().StringVarP(&valuesFile, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(buildCmd)
}
