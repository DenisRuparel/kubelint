package validators

type Severity string

const (
	Critical Severity = "CRITICAL"
	Warning  Severity = "WARNING"
	Info     Severity = "INFO"
)

type ValidationResult struct {
	File     string
	Severity Severity
	Message  string
}