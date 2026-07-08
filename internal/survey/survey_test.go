package survey

import "testing"

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		valid   bool
	}{
		{"my-app", true},
		{"my_app", true},
		{"my app", true},
		{"", false},
		{"   ", false},
		{"my/app", false},
		{"my:app", false},
		{"CON", false},
		{"nul.txt", false},
	}

	for _, tt := range tests {
		err := validateProjectName(tt.name)
		if (err == nil) != tt.valid {
			t.Fatalf("validateProjectName(%q) valid = %v, want %v", tt.name, err == nil, tt.valid)
		}
	}
}
