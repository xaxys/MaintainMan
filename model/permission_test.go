package model

import (
	"testing"
)

func TestPermissionSet(t *testing.T) {
	perms := []string{"admin.*", "user.update"}
	s := NewPermissionSet().Add(perms...)
	if !s.Has("admin.whatever") {
		t.Error("admin.whatever should be true")
	}
	if !s.Has("user.update") {
		t.Error("user.update should be true")
	}
	if s.Has("user.create") {
		t.Error("user.create should be false")
	}

	s.Add("-admin.create")
	if !s.Has("admin.whatever") {
		t.Error("admin.whatever should be true")
	}
	if s.Has("admin.create") {
		t.Error("admin.create should be false")
	}

	s.Add("-user.create")
	if s.Has("user.create") {
		t.Error("user.create should be false")
	}

	s.Add("user.create")
	if !s.Has("user.create") {
		t.Error("user.create should be true")
	}

	s.Delete("user.create")
	if s.Has("user.create") {
		t.Error("user.create should be falses")
	}

	s.Delete("user.*")
	if s.Has("user.*") {
		t.Error("user.* should be false")
	}
	if !s.Has("user.update") {
		t.Error("user.update should be true")
	}

	s.Delete("admin.*")
	if s.Has("admin.*") {
		t.Error("admin.* should be false")
	}
	if s.Has("admin.whatever") {
		t.Error("admin.whatever should be false")
	}
	if s.Has("admin.create") {
		t.Error("admin.create should be false")
	}
}
