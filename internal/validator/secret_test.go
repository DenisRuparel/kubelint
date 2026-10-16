package validators

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSecretBytes_CatchesMissingTypeAndData(t *testing.T) {
	manifest := `
apiVersion: v1
kind: Secret
metadata:
  name: empty-secret
`
	results := ValidateSecretBytes([]byte(manifest))
	if len(results) == 0 {
		t.Fatalf("expected validation findings for empty secret, got 0")
	}

	foundType := false
	foundData := false
	for _, r := range results {
		if contains(r.Message, "type field is missing") {
			foundType = true
		}
		if contains(r.Message, "must contain data or stringData") {
			foundData = true
		}
	}

	if !foundType {
		t.Errorf("expected missing type finding")
	}
	if !foundData {
		t.Errorf("expected missing data/stringData finding")
	}
}

func TestValidatePlaintextSecrets_DetectsCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	valuesPath := filepath.Join(tmpDir, "values.yaml")

	content := `
app:
  name: demo
database:
  username: dbuser
  password: admin
api:
  token: 12345
`
	if err := os.WriteFile(valuesPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test values file: %v", err)
	}

	findings := ValidatePlaintextSecrets(valuesPath)
	if len(findings) == 0 {
		t.Fatalf("expected findings for plaintext credentials in values.yaml, got 0")
	}

	foundCritical := false
	for _, f := range findings {
		if f.Severity == Critical {
			foundCritical = true
			break
		}
	}

	if !foundCritical {
		t.Errorf("expected at least one Critical finding for weak demo credentials ('admin' / '12345')")
	}
}

