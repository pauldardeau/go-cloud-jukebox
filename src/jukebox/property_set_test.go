package jukebox

import "testing"

func TestNewPropertySet(t *testing.T) {
	ps := NewPropertySet()
	if ps == nil {
		t.Log("NewPropertySet should not return nil")
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
}

func TestClear(t *testing.T) {
}

func TestContains(t *testing.T) {
}

func TestGetKeys(t *testing.T) {
}

func TestGet(t *testing.T) {
}

func TestGetIntValue(t *testing.T) {
}

func TestGetLongValue(t *testing.T) {
}

func TestGetUlongValue(t *testing.T) {
}

func TestGetBoolValue(t *testing.T) {
}

func TestGetStringValue(t *testing.T) {
}

func TestWriteToFile(t *testing.T) {
}

func TestReadFromFile(t *testing.T) {
}

func TestCount(t *testing.T) {
	ps := NewPropertySet()
	count := ps.Count()
	if count != 0 {
		t.Log("Count should return 0 on newly constructed PropertySet")
		t.Fail()
	}

	ps.Add("myInt", NewIntPropertyValue(5))
	ps.Add("myLong", NewLongPropertyValue(100))
	ps.Add("myUlong", NewUlongPropertyValue(200))
	ps.Add("myBool", NewBoolPropertyValue(true))
	ps.Add("myString", NewStringPropertyValue("foo"))
	count = ps.Count()
	if count != 5 {
		t.Fail()
	}
}

func TestToString(t *testing.T) {
}
