package enum

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Old method role enum for reference.
type OldMethodRole int

const (
	OldMethodRoleUnknown OldMethodRole = iota
	OldMethodRoleAdmin
	OldMethodRoleUser
	OldMethodRoleGuest
)

// This would be auto-generated by the enumer library.
func (r OldMethodRole) String() string {
	switch r {
	case OldMethodRoleAdmin:
		return "admin"
	case OldMethodRoleUser:
		return "user"
	case OldMethodRoleGuest:
		return "guest"
	default:
		return "unknown"
	}
}

// Old method permission enum for referece.
type OldMethodPermission int

const (
	OldMethodPermissionUnknown OldMethodPermission = iota
	OldMethodPermissionRead
	OldMethodPermissionWrite
)

// This would be auto-generated by the enumer library.
func (p OldMethodPermission) String() string {
	switch p {
	case OldMethodPermissionRead:
		return "read"
	case OldMethodPermissionWrite:
		return "write"
	default:
		return "unknown"
	}
}

// New method role enum.
type Role int

// Just to allow cleaner references.
//
// Lets assume the Role type above is defined in a package called "accounts".
// Without the type below, functions that would take an enum of type Role would
// need to be written like:
//
// func DoSomethingWithRole(r enum.Role[accounts.Role]) {}
//
// With the type below, the function can be written like:
//
// func DoSomethingWithRole(accounts.RoleEnum) {}
type RoleEnum Enum[Role]

var (
	UnknownRole = RoleEnum(New[Role]("Unknown")) // 0
	Admin       = RoleEnum(New[Role]("Admin"))   // 1
	User        = RoleEnum(New[Role]("User"))    // 2
	Guest       = RoleEnum(New[Role]("Guest"))   // 3
)

// New method permission enum.
type Permission int
type PermissionEnum Enum[Permission] // Just to allow cleaner references.

var (
	UnknownPermission = PermissionEnum(New[Permission]("Unknown")) // 0
	Read              = PermissionEnum(New[Permission]("Read"))    // 1
	Write             = PermissionEnum(New[Permission]("Write"))   // 2
)

func acceptsRoleOnly(t *testing.T, role RoleEnum) {
	t.Log(role)
}

func acceptsRoleIDOnly(t *testing.T, id Role) {
	t.Log(id)
}

func acceptsPermissionOnly(t *testing.T, permission PermissionEnum) {
	t.Log(permission)
}

func acceptsPermissionIDOnly(t *testing.T, id Permission) {
	t.Log(id)
}

func TestEnum(t *testing.T) {
	acceptsRoleOnly(t, UnknownRole)
	acceptsRoleOnly(t, Admin)
	acceptsRoleOnly(t, User)
	acceptsRoleOnly(t, Guest)

	// acceptsRoleOnly(t, UnknownPermission) // compile error

	acceptsRoleIDOnly(t, UnknownRole.ID())
	acceptsRoleIDOnly(t, Admin.ID())
	acceptsRoleIDOnly(t, User.ID())
	acceptsRoleIDOnly(t, Guest.ID())

	// acceptsRoleIDOnly(t, UnknownPermission.ID()) // compile error

	acceptsPermissionOnly(t, UnknownPermission)
	acceptsPermissionOnly(t, Read)
	acceptsPermissionOnly(t, Write)

	// acceptsPermissionOnly(t, UnknownRole) // compile error

	acceptsPermissionIDOnly(t, UnknownPermission.ID())
	acceptsPermissionIDOnly(t, Read.ID())
	acceptsPermissionIDOnly(t, Write.ID())

	// acceptsPermissionIDOnly(t, UnknownRole.ID()) // compile error
}

func TestEnum_Overflow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got normal execution")
		}
	}()

	type int8Enum int8

	// We can only have 128 int8 enums.
	for i := 0; i <= 128; i++ {
		New[int8Enum](fmt.Sprintf("Enum%d", i))
	}
}

func TestEnum_MarshalUnmarshal(t *testing.T) {
	data, err := json.Marshal(Guest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var newGuest RoleEnum
	err = json.Unmarshal(data, &newGuest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if newGuest != Guest {
		t.Errorf("expected internalEnum pointer %p, got %p", Guest.internalEnum, newGuest.internalEnum)
	}
	if newGuest.ID() != Guest.ID() {
		t.Errorf("expected ID %d, got %d", Guest.ID(), newGuest.ID())
	}
	if newGuest.String() != Guest.String() {
		t.Errorf("expected String %s, got %s", Guest.String(), newGuest.String())
	}
}

func TestEnum_Switch(t *testing.T) {
	// Unsing role values, which should be the common case.
	role := Admin

	switch role {
	case UnknownRole:
		t.Errorf("expected %s, got %s", role, UnknownRole)
	case Admin:
		// Just do not error out. This is what we want.
	case User:
		t.Errorf("expected %s, got %s", role, User)
	case Guest:
		t.Errorf("expected %s, got %s", role, Guest)
	default:
		t.Errorf("expected %s, got something else", role)
	}

	// Using IDs.

	switch roleID := role.ID(); roleID {
	case Admin.ID():
		// Just do not error out. This is what we want.
	default:
		t.Errorf("expected %d, got %d", Admin.ID(), roleID)
	}
}

func TestEnum_EnumsForType(t *testing.T) {
	enums := EnumsByType[Role]()
	if len(enums) != 4 {
		t.Errorf("expected 4, got %d", len(enums))
	}
}
