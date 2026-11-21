package models

import "testing"

func TestNewAllowedRoles(t *testing.T) {
	allowed, err := InitRoleModel(nil).NewAllowedRoles("admin", "user")
	if err != nil {
		t.Fatalf("NewAllowedRoles(admin, user) returned error: %v", err)
	}
	if !allowed.IsAllowed("admin") || !allowed.IsAllowed("user") {
		t.Fatalf("expected admin and user to be allowed: %+v", allowed)
	}

	// Unknown roles are rejected, and nothing is allowed when any role is bad.
	allowed, err = InitRoleModel(nil).NewAllowedRoles("admin", "superadmin")
	if err == nil {
		t.Fatal("NewAllowedRoles(admin, superadmin) expected an error")
	}
	if allowed.IsAllowed("admin") {
		t.Fatal("rejected role list must allow nothing, but admin was allowed")
	}

	// Role names are case-sensitive: only the lowercase system roles exist.
	if _, err := InitRoleModel(nil).NewAllowedRoles("Admin"); err == nil {
		t.Fatal("NewAllowedRoles(Admin) expected an error")
	}

	// No roles is valid, and allows nothing.
	allowed, err = InitRoleModel(nil).NewAllowedRoles()
	if err != nil {
		t.Fatalf("NewAllowedRoles() returned error: %v", err)
	}
	if allowed.IsAllowed("admin") {
		t.Fatal("empty role list must allow nothing, but admin was allowed")
	}

	// Duplicates are tolerated.
	allowed, err = InitRoleModel(nil).NewAllowedRoles("admin", "admin")
	if err != nil {
		t.Fatalf("NewAllowedRoles(admin, admin) returned error: %v", err)
	}
	if !allowed.IsAllowed("admin") {
		t.Fatal("expected admin to be allowed after duplicate entries")
	}
}
