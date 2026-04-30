package cmd

import (
	"fmt"
	"github.com/DenisRuparel/kubelint/internal/loader"
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

		if valuesFile == "" {
			valuesFile = filepath.Join(templateDir, "values.yaml")
		}

		fmt.Println("Starting KubeLint build...")
		fmt.Println("Project:", projectPath)
		fmt.Println("Templates:", templateDir)
		fmt.Println("Values file:", valuesFile)

		// Check paths
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			fmt.Println("❌ templates folder not found")
			return
		}

		if _, err := os.Stat(valuesFile); os.IsNotExist(err) {
			fmt.Println("❌ values.yaml not found")
			return
		}

		fmt.Println("✔ Paths validated successfully")

		values, err := loader.LoadValues(valuesFile)
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		fmt.Println("✔ Values loaded successfully")

		// Debug (temporary)
		fmt.Printf("Loaded values: %+v\n", values)
	},
}

func init() {
	buildCmd.Flags().StringVarP(&valuesFile, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(buildCmd)
}
