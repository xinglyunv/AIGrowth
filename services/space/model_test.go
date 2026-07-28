package space

import "testing"

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		name, role       string
		allowOwner, want bool
	}{
		{"owner allowed", RoleOwner, true, true},
		{"owner rejected for invite", RoleOwner, false, false},
		{"admin", RoleAdmin, false, true},
		{"member", RoleMember, false, true},
		{"viewer", RoleViewer, false, true},
		{"unknown", "superadmin", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRole(tt.role, tt.allowOwner); got != tt.want {
				t.Fatalf("IsValidRole(%q, %t) = %t, want %t", tt.role, tt.allowOwner, got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	for input, want := range map[string]string{"Acme Growth": "acme-growth", "  Brand@One  ": "brand-one", "!!!": "workspace"} {
		if got := slugify(input); got != want {
			t.Errorf("slugify(%q) = %q, want %q", input, got, want)
		}
	}
}
