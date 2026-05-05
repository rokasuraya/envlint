// Package audit records envlint validation runs to a newline-delimited JSON
// log file (JSON Lines format). Each entry captures the env file path, schema
// file path, UTC timestamp, all reported issues, and a top-level Passed flag.
//
// Typical usage:
//
//	entry := audit.New(envPath, schemaPath, issues)
//	if err := audit.Append("/var/log/envlint/audit.log", entry); err != nil {
//		log.Printf("audit write failed: %v", err)
//	}
//
// The log can be read back for reporting or CI history:
//
//	entries, err := audit.ReadAll("/var/log/envlint/audit.log")
//
// The file is safe to append to from multiple sequential runs; concurrent
// writes are not coordinated and should be avoided.
package audit
