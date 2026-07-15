package cli

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected [3]int
	}{
		{"0.4.2", [3]int{0, 4, 2}},
		{"v0.4.2", [3]int{0, 4, 2}},
		{"1.0.0", [3]int{1, 0, 0}},
		{"v10.20.30", [3]int{10, 20, 30}},
		{"0.10.0", [3]int{0, 10, 0}},
		{"2.0", [3]int{2, 0, 0}},
		{"v1", [3]int{1, 0, 0}},
		{"vv0.5.2", [3]int{0, 0, 0}}, // double v prefix is invalid for parseVersion but handled by TrimPrefix in caller
		{"invalid", [3]int{0, 0, 0}},
		{"", [3]int{0, 0, 0}},
	}

	for _, tt := range tests {
		result := parseVersion(tt.input)
		if result != tt.expected {
			t.Errorf("parseVersion(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		current  string
		latest   string
		expected bool // true = update available (current < latest)
	}{
		// Basic cases
		{"0.4.2", "0.4.3", true},
		{"0.4.2", "0.5.0", true},
		{"0.4.2", "1.0.0", true},

		// Same version → no update
		{"0.4.2", "0.4.2", false},
		{"v0.4.2", "v0.4.2", false},

		// Current newer → no update
		{"0.5.0", "0.4.2", false},
		{"1.0.0", "0.9.9", false},

		// THE BUG: string comparison would get these wrong
		{"0.9.0", "0.10.0", true}, // "0.9.0" > "0.10.0" as strings!
		{"0.9.9", "0.10.0", true},
		{"1.9.0", "1.10.0", true},
		{"0.99.0", "0.100.0", true},

		// With v prefix
		{"v0.4.2", "v0.5.0", true},
		{"v0.9.0", "v0.10.0", true},
		{"v1.0.0", "v0.9.0", false},

		// Mixed prefix
		{"0.4.2", "v0.5.0", true},
		{"v0.4.2", "0.5.0", true},

		// Double v prefix resolved by TrimPrefix in caller, but compareVersions directly will fail if passed raw:
		{"vv0.5.2", "v0.5.2", true}, // because vv0.5.2 parses to 0.0.0, which is < 0.5.2 (this verifies raw behavior, we fixed it in caller)
	}

	for _, tt := range tests {
		result := compareVersions(tt.current, tt.latest)
		if result != tt.expected {
			t.Errorf("compareVersions(%q, %q) = %v, want %v",
				tt.current, tt.latest, result, tt.expected)
		}
	}
}
