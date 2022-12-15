package jukebox

import "testing"

type TestHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t: t,
	}
}

func (th *TestHelper) Require(expr bool, failMessage string) {
	if !expr {
		if len(failMessage) > 0 {
			th.t.Log(failMessage)
		}
		th.t.Fail()
	}
}

func (th *TestHelper) RequireFalse(expr bool, failMessage string) {
	if expr {
		if len(failMessage) > 0 {
			th.t.Log(failMessage)
		}
		th.t.Fail()
	}
}

func (th *TestHelper) RequireStringEquals(s1 string, s2 string, failMessage string) {
	if s1 != s2 {
		if len(failMessage) > 0 {
			th.t.Log(failMessage)
		}
		th.t.Fail()
	}
}
