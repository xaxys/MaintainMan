package rbac

import (
	"testing"
)

func TestPermissionSet(t *testing.T) {
	perms := []string{"admin.*", "user.update", "tag.view.2"}
	s := newPermSet().Add(perms...)
	if !s.Has("admin.whatever") {
		t.Error("admin.whatever should be true")
	}
	if !s.Has("user.update") {
		t.Error("user.update should be true")
	}
	if s.Has("user.create") {
		t.Error("user.create should be false")
	}
	if !s.Has("tag.view.1") {
		t.Error("tag.view.1 should be true")
	}
	if !s.Has("tag.view.2") {
		t.Error("tag.view.2 should be true")
	}
	if s.Has("tag.view.3") {
		t.Error("tag.view.3 should be false")
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

	s.Add("tag.view.3")
	if !s.Has("tag.view.2") {
		t.Error("tag.view.2 should be true")
	}
	if !s.Has("tag.view.3") {
		t.Error("tag.view.3 should be true")
	}
	if s.Has("tag.view.4") {
		t.Error("tag.view.4 should be false")
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
