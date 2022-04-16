package config

import "testing"

func TestVersionCompare(t *testing.T) {
	if VersionCompare("1.2.3", "1.2.3") != 0 {
		t.Error("Expect 1.2.3 == 1.2.3")
	}
	if VersionCompare("1.2.3", "1.2.3 special") != 0 {
		t.Error("Expect 1.2.3 == 1.2.3 special")
	}
	if VersionCompare("1.2.3", "1.2.3-game") != 0 {
		t.Error("Expect 1.2.3 == 1.2.3-game")
	}
	if VersionCompare("1.2.3", "1.2.3 built-by xaxys") != 0 {
		t.Error("Expect 1.2.3 == 1.2.3 built-by xaxys")
	}

	if VersionCompare("1.2.3", "1.2.4") != -1 {
		t.Error("Expect 1.2.3 < 1.2.4")
	}
	if VersionCompare("1.2.3", "1.3.0") != -1 {
		t.Error("Expect 1.2.3 < 1.3.0")
	}
	if VersionCompare("1.2.3", "1.3") != -1 {
		t.Error("Expect 1.2.3 < 1.3")
	}
	if VersionCompare("1.2.3", "2.2.3") != -1 {
		t.Error("Expect 1.2.3 < 2.2.3")
	}
	if VersionCompare("1.2.3", "2.2.0") != -1 {
		t.Error("Expect 1.2.3 < 2.2.0")
	}
	if VersionCompare("1.2.3", "2.2") != -1 {
		t.Error("Expect 1.2.3 < 2.2")
	}
	if VersionCompare("1.2.3", "2.0") != -1 {
		t.Error("Expect 1.2.3 < 2.0")
	}
	if VersionCompare("1.2.3", "2") != -1 {
		t.Error("Expect 1.2.3 < 2")
	}

	if VersionCompare("1.2.3", "1.2.2") != 1 {
		t.Error("Expect 1.2.3 > 1.2.2")
	}
	if VersionCompare("1.2.3", "1.2.0") != 1 {
		t.Error("Expect 1.2.3 > 1.2.0")
	}
	if VersionCompare("1.2.3", "1.2") != 1 {
		t.Error("Expect 1.2.3 > 1.2")
	}
	if VersionCompare("1.2.3", "0.2.3") != 1 {
		t.Error("Expect 1.2.3 > 0.2.3")
	}
	if VersionCompare("1.2.3", "0.2.0") != 1 {
		t.Error("Expect 1.2.3 > 0.2.0")
	}
	if VersionCompare("1.2.3", "0.2") != 1 {
		t.Error("Expect 1.2.3 > 0.2")
	}
	if VersionCompare("1.2.3", "0") != 1 {
		t.Error("Expect 1.2.3 > 0")
	}
	if VersionCompare("1.2.3", "") != 1 {
		t.Error("Expect 1.2.3 > \"\"")
	}
}
