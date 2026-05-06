// Package maskwriter provides an io.Writer decorator that intercepts byte
// writes and replaces known sensitive literal values — such as passwords,
// API keys, or tokens — with the placeholder string [REDACTED].
//
// Typical usage:
//
//	sensitiveVars := map[string]string{
//		"DB_PASSWORD": "hunter2",
//		"API_KEY":     "abc123",
//	}
//	secrets := make([]string, 0, len(sensitiveVars))
//	for _, v := range sensitiveVars {
//		secrets = append(secrets, v)
//	}
//	safe := maskwriter.New(os.Stdout, secrets)
//	fmt.Fprintln(safe, "connecting with password=hunter2")
//	// Output: connecting with password=[REDACTED]
//
// The Writer is designed to be composed with any io.Writer, including
// log.New, bufio.Writer, or the formatter output writers used by envlint.
package maskwriter
