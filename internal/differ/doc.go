// Package differ computes the difference between two parsed .env maps.
//
// It is useful for detecting configuration drift between environments
// (e.g. staging vs production) or between a committed .env.example and
// a developer's local .env file.
//
// Basic usage:
//
//	baseline, _ := parser.ParseFile(".env.example")
//	target, _   := parser.ParseFile(".env")
//
//	changes := differ.Diff(baseline, target)
//	differ.WriteText(os.Stdout, changes)
//
// Sensitive values should be redacted via the redactor package before
// passing maps to WriteText to avoid leaking secrets in output.
package differ
