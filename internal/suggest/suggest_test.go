package suggest_test

import (
	"testing"

	"github.com/yourorg/envlint/internal/suggest"
)

var schemaKeys = []string{
	"DATABASE_URL",
	"DATABASE_HOST",
	"DATABASE_PORT",
	"REDIS_URL",
	"SECRET_KEY",
	"APP_ENV",
	"LOG_LEVEL",
}

func TestClosest_ExactTypo(t *testing.T) {
	got := suggest.Closest("DATABSE_URL", schemaKeys, 3)
	if len(got) == 0 {
		t.Fatal("expected at least one suggestion for 'DATABSE_URL'")
	}
	if got[0] != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL as top suggestion, got %q", got[0])
	}
}

func TestClosest_CaseInsensitive(t *testing.T) {
	got := suggest.Closest("database_url", schemaKeys, 1)
	if len(got) == 0 {
		t.Fatal("expected suggestion for lowercase input")
	}
	if got[0] != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %q", got[0])
	}
}

func TestClosest_NoMatch(t *testing.T) {
	got := suggest.Closest("ZZZZZZZZZ", schemaKeys, 3)
	if len(got) != 0 {
		t.Errorf("expected no suggestions for completely unrelated key, got %v", got)
	}
}

func TestClosest_LimitRespected(t *testing.T) {
	// DATABASE_URL and DATABASE_HOST are both close to DATABASE_UR
	got := suggest.Closest("DATABASE_UR", schemaKeys, 1)
	if len(got) != 1 {
		t.Errorf("expected exactly 1 suggestion, got %d", len(got))
	}
}

func TestClosest_OrderedByDistance(t *testing.T) {
	got := suggest.Closest("APP_EN", schemaKeys, 3)
	// APP_ENV (dist 1) should come before anything else
	if len(got) == 0 {
		t.Fatal("expected suggestions for APP_EN")
	}
	if got[0] != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %q", got[0])
	}
}

func TestClosest_EmptySchema(t *testing.T) {
	got := suggest.Closest("DATABASE_URL", []string{}, 5)
	if len(got) != 0 {
		t.Errorf("expected no suggestions for empty schema, got %v", got)
	}
}

func TestClosest_ZeroLimit(t *testing.T) {
	got := suggest.Closest("DATABASE_URL", schemaKeys, 0)
	// limit 0 means return all matches
	if len(got) == 0 {
		t.Error("expected suggestions when limit is 0 (all matches)")
	}
}
