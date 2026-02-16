package stacks

import (
	"strings"
	"testing"
)

func TestValidateSelectionRejectsUnknownFeatureWithHint(t *testing.T) {
	t.Parallel()

	sel := Selection{
		CategoryFrontend: {"frontend-htm"},
		CategoryStyling:  {"styling-tailwind"},
		CategoryHTTP:     {"http-standard"},
		CategoryDatabase: {"database-none"},
	}

	err := ValidateSelection(sel)
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "did you mean \"frontend-htmx\"") {
		t.Fatalf("expected suggestion in error, got: %s", msg)
	}
}

func TestValidateSelectionRejectsAuthWithoutSQLite(t *testing.T) {
	t.Parallel()

	sel := Selection{
		CategoryFrontend: {"frontend-htmx"},
		CategoryStyling:  {"styling-tailwind"},
		CategoryHTTP:     {"http-standard"},
		CategoryDatabase: {"database-none"},
		CategoryAuth:     {"auth-github-oauth2"},
	}

	err := ValidateSelection(sel)
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "requires \"database-sqlite\"") {
		t.Fatalf("expected dependency error, got: %s", msg)
	}
}

func TestValidateSelectionRejectsMagicLinkWithoutSQLite(t *testing.T) {
	t.Parallel()

	sel := Selection{
		CategoryFrontend: {"frontend-htmx"},
		CategoryStyling:  {"styling-tailwind"},
		CategoryHTTP:     {"http-standard"},
		CategoryDatabase: {"database-none"},
		CategoryAuth:     {"auth-magic-link"},
	}

	err := ValidateSelection(sel)
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "requires \"database-sqlite\"") {
		t.Fatalf("expected dependency error, got: %s", msg)
	}
}

func TestComposeAllowsAuthWithSQLite(t *testing.T) {
	t.Parallel()

	sel := Selection{
		CategoryFrontend: {"frontend-htmx"},
		CategoryStyling:  {"styling-tailwind"},
		CategoryHTTP:     {"http-standard"},
		CategoryDatabase: {"database-sqlite"},
		CategoryAuth:     {"auth-github-oauth2"},
	}

	stack, err := Compose(sel)
	if err != nil {
		t.Fatalf("expected compose to succeed, got: %v", err)
	}

	if !stack.HasFeature("auth-github-oauth2") {
		t.Fatal("expected composed stack to include auth-github-oauth2")
	}
}

func TestComposeAllowsMagicLinkWithSQLite(t *testing.T) {
	t.Parallel()

	sel := Selection{
		CategoryFrontend: {"frontend-htmx"},
		CategoryStyling:  {"styling-tailwind"},
		CategoryHTTP:     {"http-standard"},
		CategoryDatabase: {"database-sqlite"},
		CategoryAuth:     {"auth-magic-link"},
	}

	stack, err := Compose(sel)
	if err != nil {
		t.Fatalf("expected compose to succeed, got: %v", err)
	}

	if !stack.HasFeature("auth-magic-link") {
		t.Fatal("expected composed stack to include auth-magic-link")
	}
}
