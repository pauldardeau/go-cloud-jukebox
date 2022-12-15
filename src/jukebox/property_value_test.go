package jukebox

import "testing"

func TestNewIntPropertyValue(t *testing.T) {
	th := NewTestHelper(t)
	pv := NewIntPropertyValue(5)
	th.Require(pv != nil, "NewIntPropertyValue should not return nil")
	th.RequireStringEquals(pv.dataType, pvTypeInt, "dataType should be Int")
	th.Require(pv.intValue == 5, "intValue should match value passed to New function")
	th.Require(pv.longValue == 0, "longValue should be 0")
	th.Require(pv.ulongValue == 0, "ulongValue should be 0")
	th.RequireFalse(pv.boolValue, "boolValue should be false")
	th.RequireStringEquals(pv.stringValue, "", "stringValue should be empty string")
}

func TestNewLongPropertyValue(t *testing.T) {
	th := NewTestHelper(t)
	pv := NewLongPropertyValue(6)
	th.Require(pv != nil, "NewLongPropertyValue should not return nil")
	th.RequireStringEquals(pv.dataType, pvTypeLong, "dataType should be Long")
	th.Require(pv.intValue == 0, "intValue should be 0")
	th.Require(pv.longValue == 6, "longValue should match value passed to New function")
	th.Require(pv.ulongValue == 0, "ulongValue should be 0")
	th.RequireFalse(pv.boolValue, "boolValue should be false")
	th.RequireStringEquals(pv.stringValue, "", "stringValue should be empty string")
}

func TestNewUlongPropertyValue(t *testing.T) {
	th := NewTestHelper(t)
	pv := NewUlongPropertyValue(7)
	th.Require(pv != nil, "NewUlongPropertyValue should not return nil")
	th.RequireStringEquals(pv.dataType, pvTypeUlong, "dataType should be Ulong")
	th.Require(pv.intValue == 0, "intValue should be 0")
	th.Require(pv.longValue == 0, "longValue should be 0")
	th.Require(pv.ulongValue == 7, "ulongValue should match value passed to New function")
	th.RequireFalse(pv.boolValue, "boolValue should be false")
	th.RequireStringEquals(pv.stringValue, "", "stringValue should be empty string")
}

func TestNewBoolPropertyValue(t *testing.T) {
	th := NewTestHelper(t)
	pvTrue := NewBoolPropertyValue(true)
	th.Require(pvTrue != nil, "NewBoolPropertyValue should not return nil")
	th.RequireStringEquals(pvTrue.dataType, pvTypeBool, "dataType should be Bool")
	th.Require(pvTrue.intValue == 0, "intValue should be 0")
	th.Require(pvTrue.longValue == 0, "longValue should be 0")
	th.Require(pvTrue.ulongValue == 0, "ulongValue should be 0")
	th.Require(pvTrue.boolValue, "boolValue should match value passed to New function")
	th.RequireStringEquals(pvTrue.stringValue, "", "stringValue should be empty string")

	pvFalse := NewBoolPropertyValue(false)
	th.Require(pvFalse != nil, "NewBoolPropertyValue should not return nil")
	th.RequireStringEquals(pvFalse.dataType, pvTypeBool, "dataType should be Bool")
	th.Require(pvFalse.intValue == 0, "intValue should be 0")
	th.Require(pvFalse.longValue == 0, "longValue should be 0")
	th.Require(pvFalse.ulongValue == 0, "ulongValue should be 0")
	th.RequireFalse(pvFalse.boolValue, "boolValue should match value passed to New function")
	th.RequireStringEquals(pvFalse.stringValue, "", "stringValue should be empty string")
}

func TestNewStringPropertyValue(t *testing.T) {
	th := NewTestHelper(t)
	pv := NewStringPropertyValue("foo")
	th.Require(pv != nil, "NewStringPropertyValue should not return nil")
	th.RequireStringEquals(pv.dataType, pvTypeString, "dataType should be String")
	th.Require(pv.intValue == 0, "intValue should be 0")
	th.Require(pv.longValue == 0, "longValue should be 0")
	th.Require(pv.ulongValue == 0, "ulongValue should be 0")
	th.RequireFalse(pv.boolValue, "boolValue should be false")
	th.RequireStringEquals(pv.stringValue, "foo", "stringValue should match value passed to New function")
}

func TestGetIntValue(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.Require(pvInt.GetIntValue() == 4, "GetIntValue should match value passed to New function")
	pvLong := NewLongPropertyValue(5)
	th.Require(pvLong.GetIntValue() == 0, "GetIntValue should return 0 for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.Require(pvUlong.GetIntValue() == 0, "GetIntValue should return 0 for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.Require(pvBool.GetIntValue() == 0, "GetIntValue should return 0 for different data type")
	pvString := NewStringPropertyValue("foo")
	th.Require(pvString.GetIntValue() == 0, "GetIntValue should return 0 for different data type")
}

func TestGetLongValue(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.Require(pvInt.GetLongValue() == 0, "GetLongValue should return 0 for different data type")
	pvLong := NewLongPropertyValue(5)
	th.Require(pvLong.GetLongValue() == 5, "GetLongValue should match value passed to New function")
	pvUlong := NewUlongPropertyValue(6)
	th.Require(pvUlong.GetLongValue() == 0, "GetLongValue should return 0 for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.Require(pvBool.GetLongValue() == 0, "GetLongValue should return 0 for different data type")
	pvString := NewStringPropertyValue("foo")
	th.Require(pvString.GetLongValue() == 0, "GetLongValue should return 0 for different data type")
}

func TestGetUlongValue(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.Require(pvInt.GetUlongValue() == 0, "GetUlongValue should return 0 for different data type")
	pvLong := NewLongPropertyValue(5)
	th.Require(pvLong.GetUlongValue() == 0, "GetUlongValue should return 0 for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.Require(pvUlong.GetUlongValue() == 6, "GetUlongValue should match value passed to New function")
	pvBool := NewBoolPropertyValue(true)
	th.Require(pvBool.GetUlongValue() == 0, "GetUlongValue should return 0 for different data type")
	pvString := NewStringPropertyValue("foo")
	th.Require(pvString.GetUlongValue() == 0, "GetUlongValue should return 0 for different data type")
}

func TestGetBoolValue(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.Require(pvInt.GetBoolValue() == false, "GetBoolValue should return false for different data type")
	pvLong := NewLongPropertyValue(5)
	th.Require(pvLong.GetBoolValue() == false, "GetBoolValue should return false for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.Require(pvUlong.GetBoolValue() == false, "GetBoolValue should return false for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.Require(pvBool.GetBoolValue() == true, "GetBoolValue should match value passed to New function")
	pvString := NewStringPropertyValue("foo")
	th.Require(pvString.GetBoolValue() == false, "GetBoolValue should return false for different data type")

	pvBoolFalse := NewBoolPropertyValue(false)
	th.Require(pvBoolFalse.GetBoolValue() == false, "GetBoolValue should match value passed to New function")
}

func TestGetStringValue(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.RequireStringEquals(pvInt.GetStringValue(), "", "GetStringValue should return empty string for different data type")
	pvLong := NewLongPropertyValue(5)
	th.RequireStringEquals(pvLong.GetStringValue(), "", "GetStringValue should return empty string for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.RequireStringEquals(pvUlong.GetStringValue(), "", "GetStringValue should return empty string for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.RequireStringEquals(pvBool.GetStringValue(), "", "GetStringValue should return empty string for different data type")
	pvString := NewStringPropertyValue("foo")
	th.RequireStringEquals(pvString.GetStringValue(), "foo", "GetStringValue should match value passed to New function")
}

func TestIsInt(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.Require(pvInt.IsInt() == true, "IsInt should return true for matching data type")
	pvLong := NewLongPropertyValue(5)
	th.RequireFalse(pvLong.IsInt(), "IsInt should return false for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.RequireFalse(pvUlong.IsInt(), "IsInt should return false for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.RequireFalse(pvBool.IsInt(), "IsInt should return false for different data type")
	pvString := NewStringPropertyValue("foo")
	th.RequireFalse(pvString.IsInt(), "IsInt should return false for different data type")
}

func TestIsLong(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.RequireFalse(pvInt.IsLong(), "IsLong should return false for different data type")
	pvLong := NewLongPropertyValue(5)
	th.Require(pvLong.IsLong() == true, "IsLong should return true for matching data type")
	pvUlong := NewUlongPropertyValue(6)
	th.RequireFalse(pvUlong.IsLong(), "IsLong should return false for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.RequireFalse(pvBool.IsLong(), "IsLong should return false for different data type")
	pvString := NewStringPropertyValue("foo")
	th.RequireFalse(pvString.IsLong(), "IsLong should return false for different data type")
}

func TestIsUlong(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.RequireFalse(pvInt.IsUlong(), "IsUlong should return false for different data type")
	pvLong := NewLongPropertyValue(5)
	th.RequireFalse(pvLong.IsUlong(), "IsUlong should return false for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.Require(pvUlong.IsUlong(), "IsUlong should return true for matching data type")
	pvBool := NewBoolPropertyValue(true)
	th.RequireFalse(pvBool.IsUlong(), "IsUlong should return false for different data type")
	pvString := NewStringPropertyValue("foo")
	th.RequireFalse(pvString.IsUlong(), "IsUlong should return false for different data type")
}

func TestIsBool(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.RequireFalse(pvInt.IsBool(), "IsBool should return false for different data type")
	pvLong := NewLongPropertyValue(5)
	th.RequireFalse(pvLong.IsBool(), "IsBool should return false for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.RequireFalse(pvUlong.IsBool(), "IsBool should return false for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.Require(pvBool.IsBool(), "IsBool should return true for matching data type")
	pvString := NewStringPropertyValue("foo")
	th.RequireFalse(pvString.IsBool(), "IsBool should return false for different data type")

	pvBoolFalse := NewBoolPropertyValue(false)
	th.Require(pvBoolFalse.IsBool(), "IsBool should return true for matching data type")
}

func TestIsString(t *testing.T) {
	th := NewTestHelper(t)
	pvInt := NewIntPropertyValue(4)
	th.RequireFalse(pvInt.IsString(), "IsString should return false for different data type")
	pvLong := NewLongPropertyValue(5)
	th.RequireFalse(pvLong.IsString(), "IsString should return false for different data type")
	pvUlong := NewUlongPropertyValue(6)
	th.RequireFalse(pvUlong.IsString(), "IsString should return false for different data type")
	pvBool := NewBoolPropertyValue(true)
	th.RequireFalse(pvBool.IsString(), "IsString should return false for different data type")
	pvString := NewStringPropertyValue("foo")
	th.Require(pvString.IsString(), "IsString should return true for matching data type")
}
