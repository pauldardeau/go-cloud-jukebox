package jukebox

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	psTypeBool   string = "bool"
	psTypeString string = "string"
	psTypeInt    string = "int"
	psTypeLong   string = "long"
	psTypeUlong  string = "ulong"

	psValueTrue  string = "true"
	psValueFalse string = "false"
)

type PropertySet struct {
	mapProps map[string]*PropertyValue
}

func NewPropertySet() *PropertySet {
	var ps PropertySet
	ps.mapProps = make(map[string]*PropertyValue)
	return &ps
}

func (ps *PropertySet) Add(propName string, propValue *PropertyValue) {
	if propValue != nil {
		ps.mapProps[propName] = propValue
	}
}

func (ps *PropertySet) Clear() {
	for k := range ps.mapProps {
		delete(ps.mapProps, k)
	}
}

func (ps *PropertySet) Contains(propName string) bool {
	_, exists := ps.mapProps[propName]
	return exists
}

func (ps *PropertySet) GetKeys() []string {
	var keys []string
	for k := range ps.mapProps {
		keys = append(keys, k)
	}
	return keys
}

func (ps *PropertySet) Get(propName string) *PropertyValue {
	propValue, exists := ps.mapProps[propName]
	if exists {
		return propValue
	} else {
		return nil
	}
}

func (ps *PropertySet) GetIntValue(propName string) int {
	pv := ps.Get(propName)
	if pv != nil && pv.IsInt() {
		return pv.GetIntValue()
	} else {
		return 0
	}
}

func (ps *PropertySet) GetLongValue(propName string) int64 {
	pv := ps.Get(propName)
	if pv != nil && pv.IsLong() {
		return pv.GetLongValue()
	} else {
		return 0
	}
}

func (ps *PropertySet) GetUnsignedLongValue(propName string) uint64 {
	pv := ps.Get(propName)
	if pv != nil && pv.IsUlong() {
		return pv.GetUlongValue()
	} else {
		return 0
	}
}

func (ps *PropertySet) GetBoolValue(propName string) bool {
	pv := ps.Get(propName)
	if pv != nil && pv.IsBool() {
		return pv.GetBoolValue()
	} else {
		return false
	}
}

func (ps *PropertySet) GetStringValue(propName string) string {
	pv := ps.Get(propName)
	if pv != nil && pv.IsString() {
		return pv.GetStringValue()
	} else {
		return ""
	}
}

func (ps *PropertySet) WriteToFile(filePath string) bool {
	success := false
	s := ps.ToString()
	if len(s) > 0 {
		success = FileWriteAllText(filePath, s)
	}
	return success
}

func (ps *PropertySet) ReadFromFile(filePath string) bool {
	success := false
	fileContents, err := FileReadAllText(filePath)
	if err == nil {
		if len(fileContents) > 0 {
			fileLines := strings.Split(fileContents, "\n")
			for _, fileLine := range fileLines {
				strippedLine := strings.TrimSpace(fileLine)
				if len(strippedLine) > 0 {
					fields := strings.Split(strippedLine, "|")
					if len(fields) == 3 {
						dataType := fields[0]
						propName := fields[1]
						propValue := fields[2]

						if len(dataType) > 0 && len(propName) > 0 && len(propValue) > 0 {
							if dataType == psTypeBool {
								if propValue == psValueTrue || propValue == psValueFalse {
									boolValue := propValue == psValueTrue
									ps.Add(propName, NewBoolPropertyValue(boolValue))
								} else {
									fmt.Printf("error: invalid value for type bool '%s'\n", dataType)
									fmt.Println("skipping")
								}
							} else if dataType == psTypeString {
								ps.Add(propName, NewStringPropertyValue(propValue))
							} else if dataType == psTypeInt {
								intValue, errConv := strconv.Atoi(propValue)
								if errConv == nil {
									ps.Add(propName, NewIntPropertyValue(intValue))
								} else {
									fmt.Printf("error: unable to convert property %s value (%s) to integer\n", propName, propValue)
									return false
								}
							} else if dataType == psTypeLong {
								longValue, errConv := strconv.ParseInt(propValue, 10, 64)
								if errConv == nil {
									ps.Add(propName, NewLongPropertyValue(longValue))
								} else {
									fmt.Printf("error: unable to convert property %s value (%s) to long\n", propName, propValue)
									return false
								}
							} else if dataType == psTypeUlong {
								unsignedLongValue, errConv := strconv.ParseUint(propValue, 10, 64)
								if errConv == nil {
									ps.Add(propName, NewUlongPropertyValue(unsignedLongValue))
								} else {
									fmt.Printf("error: unable to convert property %s value (%s) to unsigned long\n", propName, propValue)
									return false
								}
							} else {
								fmt.Printf("error: unrecognized data type '%s', skipping\n", dataType)
							}
						}
					}
				}
			}
			success = true
		}
	}
	return success
}

func (ps *PropertySet) Count() int {
	return len(ps.mapProps)
}

func (ps *PropertySet) ToString() string {
	propsString := ""
	for key, pv := range ps.mapProps {
		if pv.IsBool() {
			var value string
			if pv.GetBoolValue() {
				value = psValueTrue
			} else {
				value = psValueFalse
			}
			propsString += fmt.Sprintf("%s|%s|%s\n", psTypeBool, key, value)
		} else if pv.IsString() {
			propsString += fmt.Sprintf("%s|%s|%s\n", psTypeString, key, pv.GetStringValue())
		} else if pv.IsInt() {
			propsString += fmt.Sprintf("%s|%s|%d\n", psTypeInt, key, pv.GetIntValue())
		} else if pv.IsLong() {
			propsString += fmt.Sprintf("%s|%s|%d\n", psTypeLong, key, pv.GetLongValue())
		} else if pv.IsUlong() {
			propsString += fmt.Sprintf("%s|%s|%d\n", psTypeUlong, key, pv.GetUlongValue())
		}
	}

	return propsString
}
