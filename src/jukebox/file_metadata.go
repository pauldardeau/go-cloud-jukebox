package jukebox

import (
	"fmt"
	"strconv"
)

type FileMetadata struct {
	FileUid        string
	FileName       string
	OriginFileSize int64
	StoredFileSize int64
	PadCharCount   int
	FileTime       string
	Md5Hash        string
	Compressed     bool
	Encrypted      bool
	ContainerName  string
	ObjectName     string
}

func (fm *FileMetadata) Equals(other *FileMetadata) bool {
	if other == nil {
		return false
	}
	return fm.FileUid == other.FileUid &&
		fm.FileName == other.FileName &&
		fm.OriginFileSize == other.OriginFileSize &&
		fm.StoredFileSize == other.StoredFileSize &&
		fm.PadCharCount == other.PadCharCount &&
		fm.FileTime == other.FileTime &&
		fm.Md5Hash == other.Md5Hash &&
		fm.Compressed == other.Compressed &&
		fm.Encrypted == other.Encrypted &&
		fm.ContainerName == other.ContainerName &&
		fm.ObjectName == other.ObjectName
}

func NewFileMetadata() *FileMetadata {
	var fm FileMetadata
	fm.FileUid = ""
	fm.FileName = ""
	fm.OriginFileSize = 0
	fm.StoredFileSize = 0
	fm.PadCharCount = 0
	fm.FileTime = ""
	fm.Md5Hash = ""
	fm.Compressed = false
	fm.Encrypted = false
	fm.ContainerName = ""
	fm.ObjectName = ""
	return &fm
}

func (fm *FileMetadata) FromDictionary(dictionary map[string]string) {
	fm.FromDictionaryWithPrefix(dictionary, "")
}

func (fm *FileMetadata) FromDictionaryWithPrefix(dictionary map[string]string, prefix string) {
	if dictionary != nil {
		if value, isPresent := dictionary[prefix+"FileUid"]; isPresent {
			fm.FileUid = value
		}

		if value, isPresent := dictionary[prefix+"FileName"]; isPresent {
			fm.FileName = value
		}

		if value, isPresent := dictionary[prefix+"OriginFileSize"]; isPresent {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				fm.OriginFileSize = intValue
			}
		}

		if value, isPresent := dictionary[prefix+"StoredFileSize"]; isPresent {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				fm.StoredFileSize = intValue
			}
		}

		if value, isPresent := dictionary[prefix+"PadCharCount"]; isPresent {
			intValue, err := strconv.Atoi(value)
			if err == nil {
				fm.PadCharCount = intValue
			}
		}

		if value, isPresent := dictionary[prefix+"FileTime"]; isPresent {
			fm.FileTime = value
		}

		if value, isPresent := dictionary[prefix+"Md5Hash"]; isPresent {
			fm.Md5Hash = value
		}

		if value, isPresent := dictionary[prefix+"Compressed"]; isPresent {
			fm.Compressed = value == "1"
		}

		if value, isPresent := dictionary[prefix+"Encrypted"]; isPresent {
			fm.Encrypted = value == "1"
		}

		if value, isPresent := dictionary[prefix+"ContainerName"]; isPresent {
			fm.ContainerName = value
		}

		if value, isPresent := dictionary[prefix+"ObjectName"]; isPresent {
			fm.ObjectName = value
		}
	}
}

func (fm *FileMetadata) ToDictionary() map[string]string {
	return fm.ToDictionaryWithPrefix("")
}

func (fm *FileMetadata) ToDictionaryWithPrefix(prefix string) map[string]string {
	compressedValue := "0"
	encryptedValue := "0"
	if fm.Compressed {
		compressedValue = "1"
	}
	if fm.Encrypted {
		encryptedValue = "1"
	}
	return map[string]string{
		prefix + "FileUid":        fm.FileUid,
		prefix + "FileName":       fm.FileName,
		prefix + "OriginFileSize": fmt.Sprintf("%d", fm.OriginFileSize),
		prefix + "StoredFileSize": fmt.Sprintf("%d", fm.StoredFileSize),
		prefix + "PadCharCount":   fmt.Sprintf("%d", fm.PadCharCount),
		prefix + "FileTime":       fm.FileTime,
		prefix + "Md5Hash":        fm.Md5Hash,
		prefix + "Compressed":     compressedValue,
		prefix + "Encrypted":      encryptedValue,
		prefix + "ContainerName":  fm.ContainerName,
		prefix + "ObjectName":     fm.ObjectName,
	}
}

func (fm *FileMetadata) ToPropertySet() *PropertySet {
	return fm.ToPropertySetWithPrefix("")
}

func (fm *FileMetadata) ToPropertySetWithPrefix(prefix string) *PropertySet {
	props := NewPropertySet()
	props.Add(prefix+"file_uid", NewStringPropertyValue(fm.FileUid))
	props.Add(prefix+"file_name", NewStringPropertyValue(fm.FileName))
	props.Add(prefix+"origin_file_size", NewLongPropertyValue(fm.OriginFileSize))
	props.Add(prefix+"stored_file_size", NewLongPropertyValue(fm.StoredFileSize))
	props.Add(prefix+"pad_char_count", NewIntPropertyValue(fm.PadCharCount))
	props.Add(prefix+"file_time", NewStringPropertyValue(fm.FileTime))
	props.Add(prefix+"md5_hash", NewStringPropertyValue(fm.Md5Hash))
	props.Add(prefix+"compressed", NewBoolPropertyValue(fm.Compressed))
	props.Add(prefix+"encrypted", NewBoolPropertyValue(fm.Encrypted))
	props.Add(prefix+"container_name", NewStringPropertyValue(fm.ContainerName))
	props.Add(prefix+"object_name", NewStringPropertyValue(fm.ObjectName))
	return props
}
