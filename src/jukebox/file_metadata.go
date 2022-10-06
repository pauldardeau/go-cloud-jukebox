package jukebox

import (
   "fmt"
   "strconv"
)

type FileMetadata struct {
   File_uid string
   File_name string
   Origin_file_size int64
   Stored_file_size int64
   Pad_char_count int
   File_time string
   Md5_hash string
   Compressed bool
   Encrypted bool
   Container_name string
   Object_name string
}

func (fm *FileMetadata) Equals(other *FileMetadata) bool {
   if other == nil {
      return false
   }
   return fm.File_uid == other.File_uid &&
          fm.File_name == other.File_name &&
          fm.Origin_file_size == other.Origin_file_size &&
          fm.Stored_file_size == other.Stored_file_size &&
          fm.Pad_char_count == other.Pad_char_count &&
          fm.File_time == other.File_time &&
          fm.Md5_hash == other.Md5_hash &&
          fm.Compressed == other.Compressed &&
          fm.Encrypted == other.Encrypted &&
          fm.Container_name == other.Container_name &&
          fm.Object_name == other.Object_name
}

func NewFileMetadata() *FileMetadata {
   var fm FileMetadata
   fm.File_uid = ""
   fm.File_name = ""
   fm.Origin_file_size = 0
   fm.Stored_file_size = 0
   fm.Pad_char_count = 0
   fm.File_time = ""
   fm.Md5_hash = ""
   fm.Compressed = false
   fm.Encrypted = false
   fm.Container_name = ""
   fm.Object_name = ""
   return &fm
}

func (fm *FileMetadata) From_Dictionary(dictionary map[string]string) {
   fm.From_Dictionary_With_Prefix(dictionary, "")
}

func (fm *FileMetadata) From_Dictionary_With_Prefix(dictionary map[string]string, prefix string) {
   if dictionary != nil {
      if value, isPresent := dictionary[prefix + "file_uid"]; isPresent {
         fm.File_uid = value
      }

      if value, isPresent := dictionary[prefix + "file_name"]; isPresent {
         fm.File_name = value
      }

      if value, isPresent := dictionary[prefix + "origin_file_size"]; isPresent {
         int_value, err := strconv.ParseInt(value, 10, 64)
         if err == nil {
            fm.Origin_file_size = int_value
         }
      }

      if value, isPresent := dictionary[prefix + "stored_file_size"]; isPresent {
         int_value, err := strconv.ParseInt(value, 10, 64)
         if err == nil {
            fm.Stored_file_size = int_value
         }
      }

      if value, isPresent := dictionary[prefix + "pad_char_count"]; isPresent {
         int_value, err := strconv.Atoi(value)
         if err == nil {
            fm.Pad_char_count = int_value
         }
      }

      if value, isPresent := dictionary[prefix + "file_time"]; isPresent {
         fm.File_time = value
      }

      if value, isPresent := dictionary[prefix + "md5_hash"]; isPresent {
         fm.Md5_hash = value
      }

      if value, isPresent := dictionary[prefix + "compressed"]; isPresent {
         fm.Compressed = (value == "1")
      }

      if value, isPresent := dictionary[prefix + "encrypted"]; isPresent {
         fm.Encrypted = (value == "1")
      }

      if value, isPresent := dictionary[prefix + "container_name"]; isPresent {
         fm.Container_name = value
      }

      if value, isPresent := dictionary[prefix + "object_name"]; isPresent {
         fm.Object_name = value
      }
   }
}

func (fm *FileMetadata) To_Dictionary() map[string]string {
   return fm.To_Dictionary_With_Prefix("")
}

func (fm *FileMetadata) To_Dictionary_With_Prefix(prefix string) map[string]string {
   compressed_value := "0"
   encrypted_value := "0"
   if fm.Compressed {
      compressed_value = "1"
   }
   if fm.Encrypted {
      encrypted_value = "1"
   }
   return map[string]string{
      prefix + "file_uid": fm.File_uid,
      prefix + "file_name": fm.File_name,
      prefix + "origin_file_size": fmt.Sprintf("%d", fm.Origin_file_size),
      prefix + "stored_file_size": fmt.Sprintf("%d", fm.Stored_file_size),
      prefix + "pad_char_count": fmt.Sprintf("%d", fm.Pad_char_count),
      prefix + "file_time": fm.File_time,
      prefix + "md5_hash": fm.Md5_hash,
      prefix + "compressed": compressed_value,
      prefix + "encrypted": encrypted_value,
      prefix + "container_name": fm.Container_name,
      prefix + "object_name": fm.Object_name,
   }
}
