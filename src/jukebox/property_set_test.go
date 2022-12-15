package jukebox

import "testing"

func TestNewPropertySet(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	th.Require(ps != nil, "NewPropertySet should not return nil")
	th.Require(ps.mapProps != nil, "NewPropertySet should initialize mapProps")
	th.Require(ps.Count() == 0, "NewPropertySet should have no properties")
}

func TestAdd(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	ps.Add("foo", NewIntPropertyValue(5))
	ps.Add("bar", NewBoolPropertyValue(true))
	th.Require(ps.Count() == 2, "There should be 2 properties after 2 calls to Add")
}

func TestClear(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	ps.Clear()
	ps.Add("foo", NewIntPropertyValue(5))
	ps.Add("bar", NewBoolPropertyValue(true))
	ps.Clear()
	th.Require(ps.Count() == 0, "Clear should result in having no properties")
}

func TestContains(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	th.RequireFalse(ps.Contains("foo"), "Contains should return false for non-existing property")
	ps.Add("bar", NewStringPropertyValue("baz"))
	th.Require(ps.Contains("bar"), "Contains should return true for existing property")
}

func TestGetKeys(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	emptyKeys := ps.GetKeys()
	th.Require(len(emptyKeys) == 0, "GetKeys should return empty slice when no keys are present")
	ps.Add("venus", NewStringPropertyValue("Venus"))
	ps.Add("mars", NewStringPropertyValue("Mars"))
	ps.Add("mercury", NewStringPropertyValue("Mercury"))
	keys := ps.GetKeys()
	th.Require(len(keys) == 3, "GetKeys should return 3 after 3 calls of Add")
	venusCount := 0
	marsCount := 0
	mercuryCount := 0
	otherCount := 0
	for _, key := range keys {
		if key == "venus" {
			venusCount += 1
		} else if key == "mars" {
			marsCount += 1
		} else if key == "mercury" {
			mercuryCount += 1
		} else {
			otherCount += 1
		}
	}
	th.Require(venusCount == 1, "should be 1 venus in GetKeys")
	th.Require(marsCount == 1, "should be 1 mars in GetKeys")
	th.Require(mercuryCount == 1, "should be 1 mercury in GetKeys")
	th.Require(otherCount == 0, "should be no other keys than what was used in Add calls")
}

func TestGet(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	nonExistingProp := ps.Get("foo")
	th.Require(nonExistingProp == nil, "Get should return nil for non-existing property")
	ps.Add("bar", NewLongPropertyValue(100))
	pvBar := ps.Get("bar")
	th.Require(pvBar != nil, "Get should return non-nil for existing property")
	th.Require(pvBar.GetLongValue() == 100, "property value should have same value it was created with")
}

func TestPSGetIntValue(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	intValue := ps.GetIntValue("foo")
	th.Require(intValue == 0, "GetIntValue should return 0 for non-existing property")
	ps.Add("numPlanets", NewIntPropertyValue(9))
	intValue = ps.GetIntValue("numPlanets")
	th.Require(intValue == 9, "GetIntValue should return matching value used in Add")
}

func TestPSGetLongValue(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	longValue := ps.GetLongValue("foo")
	th.Require(longValue == 0, "GetLongValue should return 0 for non-existing property")
	ps.Add("numPlanets", NewLongPropertyValue(9))
	longValue = ps.GetLongValue("numPlanets")
	th.Require(longValue == 9, "GetLongValue should return matching value used in Add")
}

func TestPSGetUlongValue(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	ulongValue := ps.GetUlongValue("foo")
	th.Require(ulongValue == 0, "GetUlongValue should return 0 for non-existing property")
	ps.Add("numPlanets", NewUlongPropertyValue(9))
	ulongValue = ps.GetUlongValue("numPlanets")
	th.Require(ulongValue == 9, "GetUlongValue should return matching value used in Add")
}

func TestPSGetBoolValue(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	boolValue := ps.GetBoolValue("foo")
	th.Require(boolValue == false, "GetBoolValue should return false for non-existing property")
	ps.Add("havePlanets", NewBoolPropertyValue(true))
	boolValue = ps.GetBoolValue("havePlanets")
	th.Require(boolValue == true, "GetBoolValue should return matching value used in Add")

	ps.Add("haveNoPlanets", NewBoolPropertyValue(false))
	boolValue = ps.GetBoolValue("haveNoPlanets")
	th.Require(boolValue == false, "GetBoolValue should return matching value used in Add")
}

func TestPSGetStringValue(t *testing.T) {
	th := NewTestHelper(t)
	ps := NewPropertySet()
	stringValue := ps.GetStringValue("foo")
	th.RequireStringEquals(stringValue, "", "GetStringValue should return empty string for non-existing property")
	ps.Add("ourPlanet", NewStringPropertyValue("Earth"))
	stringValue = ps.GetStringValue("ourPlanet")
	th.RequireStringEquals(stringValue, "Earth", "GetStringValue should return matching value used in Add")
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
