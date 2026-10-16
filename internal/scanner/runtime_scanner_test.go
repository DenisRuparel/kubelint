package scanner

import (
	"strings"
	"testing"
)

func TestSplitYAMLDocuments_DoesNotSplitInsideCertificates(t *testing.T) {
	manifestWithCert := `apiVersion: v1
kind: Secret
metadata:
  name: tls-cert
type: kubernetes.io/tls
data:
  tls.crt: |
    -----BEGIN CERTIFICATE-----
    MIIBkTCB+wIBADANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZteS1jYTEUMBIG
    -----END CERTIFICATE-----
---
apiVersion: v1
kind: Service
metadata:
  name: web-service
spec:
  ports:
  - port: 80
`

	docs := splitYAMLDocuments(manifestWithCert)
	if len(docs) != 2 {
		t.Fatalf("expected 2 documents, got %d. Did it incorrectly split on '-----'?", len(docs))
	}

	if !strings.Contains(docs[0], "tls-cert") {
		t.Errorf("expected doc 0 to be tls-cert secret, got: %s", docs[0])
	}
	if !strings.Contains(docs[1], "web-service") {
		t.Errorf("expected doc 1 to be web-service, got: %s", docs[1])
	}
}

func TestScanRenderedYAML_MultiDocScan(t *testing.T) {
	manifests := `apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  ports:
  - port: 80
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: app
        image: nginx:1.25
`

	result := ScanRenderedYAML(manifests)
	// Should have scanned both documents without crashing or syntax errors
	if result.Summary.Critical > 0 {
		for _, issue := range result.Issues {
			t.Logf("issue: [%v] %s", issue.Severity, issue.Message)
		}
	}
	if len(result.Issues) == 0 {
		t.Errorf("expected at least some warnings (e.g. securityContext or probes), got 0")
	}
}

