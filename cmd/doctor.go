package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [project-path]",
	Short: "Check system and project readiness for KubeLint (defaults to current directory)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var projectPath string

		if len(args) == 0 {
			projectPath, _ = os.Getwd()
		} else {
			projectPath = args[0]
		}

		fmt.Println("\n🔍 KubeLint Doctor Report")
		fmt.Println("-----------------------------------------")
		fmt.Println()

		section("📦 Environment")
		checkCommand("kubectl", "kubectl")
		checkContainerRuntime()
		checkOptionalCommand("cue", "CUE (recommended for Phase 3)")

		fmt.Println()

		section("☸️ Kubernetes")

		checkCluster()
		checkContext()
		checkKubectlVersion()

		fmt.Println()

		section("📁 Project")

		checkProjectStructure(projectPath)

		fmt.Println("------------------------------------------")
		fmt.Println("✅ Doctor check completed")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

//////////////////////////////////////////////////////
// UI HELPERS
//////////////////////////////////////////////////////

func section(title string) {
	fmt.Println(title)
	fmt.Println("------------------------------------------")
}

func success(msg string) {
	fmt.Printf("\033[32m✔ %s\033[0m\n", msg)
}

func warn(msg string) {
	fmt.Printf("\033[33m⚠ %s\033[0m\n", msg)
}

func fail(msg string) {
	fmt.Printf("\033[31m✖ %s\033[0m\n", msg)
}

//////////////////////////////////////////////////////
// CHECKS
//////////////////////////////////////////////////////

func checkCommand(command string, name string) {
	_, err := exec.LookPath(command)
	if err != nil {
		fail(fmt.Sprintf("%s is NOT installed", name))
		return
	}
	success(fmt.Sprintf("%s is installed", name))
}

func checkContainerRuntime() {
	_, dockerErr := exec.LookPath("docker")
	_, ctrErr := exec.LookPath("containerd")

	if dockerErr == nil {
		success("Docker is installed")
		return
	}

	if ctrErr == nil {
		success("containerd is installed")
		return
	}

	fail("No container runtime found (Docker or containerd required)")
}

func checkOptionalCommand(command string, name string) {
	_, err := exec.LookPath(command)
	if err != nil {
		warn(fmt.Sprintf("%s is not installed", name))
		return
	}
	success(fmt.Sprintf("%s is installed", name))
}

//////////////////////////////////////////////////////
// KUBERNETES CHECKS
//////////////////////////////////////////////////////

func checkCluster() {
	cmd := exec.Command("kubectl", "cluster-info")
	err := cmd.Run()
	if err != nil {
		fail("Cannot connect to Kubernetes cluster")
		fmt.Println("   💡 Run: kubectl config use-context <context>")
		return
	}
	success("Kubernetes cluster reachable")
}

func checkContext() {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		fail("Unable to get current context")
		return
	}

	context := strings.TrimSpace(string(out))
	success(fmt.Sprintf("Current Context: %s", context))

	nsOut, err := exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={..namespace}").Output()
	if err != nil {
		warn("Could not determine default namespace")
		return
	}

	namespace := strings.TrimSpace(string(nsOut))
	if namespace == "" {
		namespace = "default"
	}
	success(fmt.Sprintf("Default Namespace: %s", namespace))
}

func checkKubectlVersion() {
	out, _ := exec.Command("kubectl", "version", "--client").CombinedOutput()

	if len(out) == 0 {
		fail("kubectl installed but no version output")
		return
	}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		l := strings.ToLower(line)

		if strings.Contains(l, "client version") || strings.Contains(l, "clientversion") {

			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				version := strings.TrimSpace(parts[1])
				success(fmt.Sprintf("kubectl version: %s", version))
				return
			}

			// fallback if split fails
			success(strings.TrimSpace(line))
			return
		}
	}

	// fallback: print raw output (important)
	warn("Could not parse kubectl version, raw output:")
	fmt.Println(string(out))
}

//////////////////////////////////////////////////////
// PROJECT CHECKS
//////////////////////////////////////////////////////

func checkProjectStructure(projectPath string) {

	templatesPath := filepath.Join(projectPath, "templates")
	valuesPath := filepath.Join(projectPath, "templates", "values.yaml")

	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		fail("templates/ folder not found (not a KubeLint project)")
		return
	}

	success("templates/ folder found")

	if _, err := os.Stat(valuesPath); os.IsNotExist(err) {
		fail("values.yaml not found inside templates/")
		return
	}

	success("values.yaml found")

	files, err := os.ReadDir(templatesPath)
	if err != nil || len(files) == 0 {
		warn("No template files found")
		return
	}

	success(fmt.Sprintf("Found %d template files", len(files)))
}
