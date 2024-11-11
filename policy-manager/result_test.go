package policy_manager

import "testing"

func TestViolation_GetString(t *testing.T) {
	t.Run("Fetch key from violation returns key", func(t *testing.T) {
		violation := &Violation{
			"title": "Some violation was found",
		}
		result := violation.GetString("title", "Some default")
		if result != "Some violation was found" {
			t.Errorf("GetString() = %s, expected %s", result, "Some violation was found")
		}
	})

	t.Run("Fetch non-existent key from violation returns default", func(t *testing.T) {
		violation := &Violation{}
		result := violation.GetString("title", "Some default")
		if result != "Some default" {
			t.Errorf("GetString() = %s, expected %s", result, "Some default")
		}
	})
}
