package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system readiness for KubeLint",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running KubeLint system diagnostics...")
		fmt.Println()
		checkCommand("go", "Go")
		checkCommand("kubectl", "kubectl")

		fmt.Println("\nSystem diagnostics completed.")
	},
}

func checkCommand(command string, name string) {
	_, err := exec.LookPath(command)
	if err != nil {
		fmt.Printf("✗ %s is NOT installed\n", name)
		return
	}

	fmt.Printf("✓ %s is installed\n", name)
}