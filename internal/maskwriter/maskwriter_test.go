package maskwriter_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envlint/internal/maskwriter"
)

func TestWrite_MasksSingleSecret(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, []string{"s3cr3t"})

	_, err := mw.Write([]byte("password=s3cr3t\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "password=[REDACTED]\n" {
		t.Errorf("got %q, want %q", got, "password=[REDACTED]\n")
	}
}

func TestWrite_MasksMultipleSecrets(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, []string{"alpha", "beta"})

	_, _ = mw.Write([]byte("a=alpha b=beta\n"))
	if got := buf.String(); got != "a=[REDACTED] b=[REDACTED]\n" {
		t.Errorf("got %q", got)
	}
}

func TestWrite_NoSecretsPassthrough(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, nil)

	_, _ = mw.Write([]byte("hello world\n"))
	if got := buf.String(); got != "hello world\n" {
		t.Errorf("got %q", got)
	}
}

func TestWrite_EmptySecretIgnored(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, []string{"", "real"})

	_, _ = mw.Write([]byte("val=real\n"))
	if got := buf.String(); got != "val=[REDACTED]\n" {
		t.Errorf("got %q", got)
	}
}

func TestWrite_ReturnsOriginalLength(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, []string{"secret"})

	input := []byte("key=secret")
	n, err := mw.Write(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("n = %d, want %d", n, len(input))
	}
}

func TestAddSecret_MasksAfterConstruction(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, nil)
	mw.AddSecret("late")

	_, _ = mw.Write([]byte("token=late\n"))
	if got := buf.String(); got != "token=[REDACTED]\n" {
		t.Errorf("got %q", got)
	}
}

func TestAddSecret_EmptyIgnored(t *testing.T) {
	var buf bytes.Buffer
	mw := maskwriter.New(&buf, nil)
	mw.AddSecret("")

	_, _ = mw.Write([]byte("plain\n"))
	if got := buf.String(); got != "plain\n" {
		t.Errorf("got %q", got)
	}
}
