package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/DenisRuparel/kubelint/internal/builder"
	"github.com/spf13/cobra"
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
			return
		}

		// kubectl apply -f -
		kubectl := exec.Command("kubectl", "apply", "-f", "-")
		kubectl.Stdin = strings.NewReader(output)

		out, err := kubectl.CombinedOutput()
		if err != nil {
			fmt.Println("❌ kubectl apply failed:")
			fmt.Println(string(out))
			return
		}

		fmt.Println("✔ Applied successfully")
		fmt.Println(string(out))
	},
}

func init() {
	applyCmd.Flags().StringVarP(&applyValues, "values", "f", "", "Path to values.yaml")
	rootCmd.AddCommand(applyCmd)
}