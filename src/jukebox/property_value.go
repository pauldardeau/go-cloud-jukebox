package jukebox

type PropertyValue struct {
	dataType string
	// conceptually, the following members are a union data type.
	// however, go doesn't have union data types. there are some
	// go tricks that could be employed to make an instance of
	// this type require less memory space, but I'm choosing to
	// not do that because it makes the code harder to understand.
	intValue    int
	longValue   int64
	ulongValue  uint64
	boolValue   bool
	stringValue string
}

const (
	pvTypeInt    string = "Int"
	pvTypeLong   string = "Long"
	pvTypeUlong  string = "ULong"
	pvTypeBool   string = "Bool"
	pvTypeString string = "String"
)

func NewIntPropertyValue(intValue int) *PropertyValue {
	return &PropertyValue{pvTypeInt, intValue, 0, 0, false, ""}
}

func NewLongPropertyValue(longValue int64) *PropertyValue {
	return &PropertyValue{pvTypeLong, 0, longValue, 0, false, ""}
}

func NewUlongPropertyValue(ulongValue uint64) *PropertyValue {
	return &PropertyValue{pvTypeUlong, 0, 0, ulongValue, false, ""}
}

func NewBoolPropertyValue(boolValue bool) *PropertyValue {
	return &PropertyValue{pvTypeBool, 0, 0, 0, boolValue, ""}
}

func NewStringPropertyValue(stringValue string) *PropertyValue {
	return &PropertyValue{pvTypeString, 0, 0, 0, false, stringValue}
}

func (pv *PropertyValue) GetIntValue() int {
	if pv.IsInt() {
		return pv.intValue
	} else {
		return 0
	}
}

func (pv *PropertyValue) GetLongValue() int64 {
	if pv.IsLong() {
		return pv.longValue
	} else {
		return 0
	}
}

func (pv *PropertyValue) GetUlongValue() uint64 {
	if pv.IsUlong() {
		return pv.ulongValue
	} else {
		return 0
	}
}

func (pv *PropertyValue) GetBoolValue() bool {
	if pv.IsBool() {
		return pv.boolValue
	} else {
		return false
	}
}

func (pv *PropertyValue) GetStringValue() string {
	if pv.IsString() {
		return pv.stringValue
	} else {
		return ""
	}
}

func (pv *PropertyValue) IsInt() bool {
	return pv.dataType == pvTypeInt
}

func (pv *PropertyValue) IsLong() bool {
	return pv.dataType == pvTypeLong
}

func (pv *PropertyValue) IsUlong() bool {
	return pv.dataType == pvTypeUlong
}

func (pv *PropertyValue) IsBool() bool {
	return pv.dataType == pvTypeBool
}

func (pv *PropertyValue) IsString() bool {
	return pv.dataType == pvTypeString
}
