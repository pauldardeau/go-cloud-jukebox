package jukebox

import "testing"

func TestNewTestHelper(t *testing.T) {
	th := NewTestHelper(t)
	if th == nil {
		t.Log("NewTestHelper should not return nil")
		t.Fail()
	}
	if th.t != t {
		t.Log("t should match value passed to New function")
		t.Fail()
	}
}

func TestRequire(t *testing.T) {
	th := NewTestHelper(t)
	th.Require(true, "true should not fail for Require")
}

func TestRequireFalse(t *testing.T) {
	th := NewTestHelper(t)
	th.RequireFalse(false, "false should not fail for RequireFalse")
}

func TestRequireStringEquals(t *testing.T) {
	th := NewTestHelper(t)
	th.RequireStringEquals("foo", "foo", "matching strings should not fail for RequireStringEquals")
}
