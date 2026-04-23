package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "APP_ENV" || entries[0].Value != "production" {
		t.Errorf("entry[0] = %+v, want {APP_ENV production}", entries[0])
	}
	if entries[1].Key != "PORT" || entries[1].Value != "8080" {
		t.Errorf("entry[1] = %+v, want {PORT 8080}", entries[1])
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"` + "\nSECRET='abc123'\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "postgres://localhost/mydb" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
	if entries[1].Value != "abc123" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export MY_VAR=hello\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "MY_VAR" || entries[0].Value != "hello" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseFile_MissingEquals(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for line missing '=', got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
