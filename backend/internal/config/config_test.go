package config

import "testing"

// TestLoadRequiresDatabaseURL 验证该路径的预期行为。
func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}
