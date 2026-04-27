// Package expander provides variable interpolation for .env file values.
//
// It supports two reference syntaxes:
//
//	${VAR_NAME}   — brace-delimited reference
//	$VAR_NAME     — bare reference (terminated by a non-identifier character)
//
// A literal dollar sign can be produced with $$.
//
// References are resolved first from the provided env map, then optionally
// from the OS environment when FallbackToOS is true (the default).
//
// Example:
//
//	env := map[string]string{
//		"HOST": "db.internal",
//		"DSN":  "postgres://${HOST}/mydb",
//	}
//	ex := expander.New()
//	_ = ex.Expand(env)
//	// env["DSN"] == "postgres://db.internal/mydb"
package expander
