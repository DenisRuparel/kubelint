package validators

import (
	"testing"
)

func TestValidateSecurityBytes_CatchesInsecureDeployment(t *testing.T) {
	insecureManifest := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: insecure-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: insecure-app
  template:
    metadata:
      labels:
        app: insecure-app
    spec:
      containers:
      - name: web
        image: nginx:latest
`

	results := ValidateSecurityBytes([]byte(insecureManifest))

	if len(results) == 0 {
		t.Fatalf("expected security issues to be detected for insecure deployment, but got 0 results (Bug #1 regression)")
	}

	foundRunAsNonRoot := false
	foundPrivileged := false
	foundReadOnlyRoot := false

	for _, r := range results {
		if r.Severity != Warning && r.Severity != Critical {
			t.Errorf("unexpected severity: %v", r.Severity)
		}
		if r.Message != "" {
			switch {
			case contains(r.Message, "runAsNonRoot"):
				foundRunAsNonRoot = true
			case contains(r.Message, "privileged"):
				foundPrivileged = true
			case contains(r.Message, "readOnlyRootFilesystem"):
				foundReadOnlyRoot = true
			}
		}
	}

	if !foundRunAsNonRoot {
		t.Errorf("expected warning for runAsNonRoot")
	}
	if !foundPrivileged {
		t.Errorf("expected warning for privileged")
	}
	if !foundReadOnlyRoot {
		t.Errorf("expected warning for readOnlyRootFilesystem")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && searchSubstr(s, substr)))
}

func searchSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

