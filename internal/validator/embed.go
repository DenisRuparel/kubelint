package validators

import "embed"

// 🔥 Now valid (same directory or child directory)
//go:embed schemas/*.cue
var schemaFS embed.FS