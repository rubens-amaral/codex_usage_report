package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileWithFractionalPrimaryPercent(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "rollout-test.jsonl")


	content := `{"timestamp":"2024-01-01T00:00:00Z","payload":{"info":{"total_token_usage":{"total_tokens":10},"last_token_usage":{"total_tokens":5}},"rate_limits":{"primary":{"used_percent":0.42},"secondary":{"used_percent":0.0}}}}`

	if err := os.WriteFile(filePath, []byte(content+"\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	timelineFull, _, _, _, err := ParseFile(filePath, false)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if len(timelineFull) == 0 {
		t.Fatalf("expected at least one timeline entry, got 0")
	}


	entry := timelineFull[0]
	if entry.Primary != 42 {
		t.Fatalf("expected primary percent to round to 42, got %d", entry.Primary)
	}

	if entry.Secondary != 0 {
		t.Fatalf("expected secondary percent to stay at 0, got %d", entry.Secondary)
	}
}

func TestParseFileWithLegacyFlatRateLimitFields(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "rollout-legacy.jsonl")

	content := `{"timestamp":"2024-01-02T00:00:00Z","payload":{"info":{"total_token_usage":{"total_tokens":20},"last_token_usage":{"total_tokens":10}},"rate_limits":{"primary_used_percent":15.5,"secondary_used_percent":51.2}}}`
	if err := os.WriteFile(filePath, []byte(content+"\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	timelineFull, _, _, _, err := ParseFile(filePath, false)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if len(timelineFull) != 1 {
		t.Fatalf("expected exactly one entry, got %d", len(timelineFull))
	}

	entry := timelineFull[0]
	if entry.Primary != 16 {
		t.Fatalf("expected primary to round to 16, got %d", entry.Primary)
	}

	if entry.Secondary != 51 {
		t.Fatalf("expected secondary to round to 51, got %d", entry.Secondary)
	}
}


func TestParseFileWithNestedRateLimitStructure(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "rollout-nested.jsonl")

	content := `{"timestamp":"2025-09-29T04:16:49.337Z","type":"event_msg","payload":{"type":"token_count","info":{"total_token_usage":{"total_tokens":10},"last_token_usage":{"total_tokens":5}},"rate_limits":{"primary":{"used_percent":6.0},"secondary":{"used_percent":100.0}}}}`
	if err := os.WriteFile(filePath, []byte(content+"\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	timelineFull, _, _, _, err := ParseFile(filePath, false)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if len(timelineFull) != 1 {
		t.Fatalf("expected exactly one timeline entry, got %d", len(timelineFull))
	}

	entry := timelineFull[0]
	if entry.Primary != 6 {
		t.Fatalf("expected primary to be 6, got %d", entry.Primary)
	}

	if entry.Secondary != 100 {
		t.Fatalf("expected secondary to be 100, got %d", entry.Secondary)
	}
}

