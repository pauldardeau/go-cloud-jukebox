package jukebox

import (
   "fmt"
)

type JukeboxOptions struct {
   DebugMode bool
   UseEncryption bool
   UseCompression bool
   CheckDataIntegrity bool
   FileCacheCount int
   NumberSongs int
   EncryptionKey string
   EncryptionKeyFile string
   EncryptionIv string
   SuppressMetadataDownload bool
}

func NewJukeboxOptions() (*JukeboxOptions) {
   var o JukeboxOptions
   o.DebugMode = false
   o.UseEncryption = false
   o.UseCompression = false
   o.CheckDataIntegrity = false
   o.FileCacheCount = 3
   o.NumberSongs = 0
   o.EncryptionKey = ""
   o.EncryptionKeyFile = ""
   o.EncryptionIv = ""
   o.SuppressMetadataDownload = false
   return &o
}

func printBoolValue(varName string, boolValue bool) {
   if boolValue {
      fmt.Printf("%s = true\n", varName)
   } else {
      fmt.Printf("%s = false\n", varName)
   }
}

func (o *JukeboxOptions) Show() {
   fmt.Println("========= Start JukeboxOptions ========")
   printBoolValue("DebugMode", o.DebugMode)
   printBoolValue("UseEncryption", o.UseEncryption)
   printBoolValue("UseCompression", o.UseCompression)
   printBoolValue("CheckDataIntegrity", o.CheckDataIntegrity)
   fmt.Printf("FileCacheCount = %d\n", o.FileCacheCount)
   fmt.Printf("NumberSongs = %d\n", o.NumberSongs)
   fmt.Printf("EncryptionKey = '%s'\n", o.EncryptionKey)
   fmt.Printf("EncryptionKeyFile = '%s'\n", o.EncryptionKeyFile)
   fmt.Printf("EncryptionIv = '%s'\n", o.EncryptionIv)
   printBoolValue("SuppressMetadataDownload", o.SuppressMetadataDownload)
   fmt.Println("========= End JukeboxOptions =========")
}

func (o *JukeboxOptions) ValidateOptions() (bool) {
   if o.FileCacheCount < 0 {
      fmt.Println("error: file cache count must be non-negative integer value")
      return false
   }

   if len(o.EncryptionKeyFile) > 0 && ! FileExists(o.EncryptionKeyFile) {
      fmt.Printf("error: encryption key file doesn't exist '%s'\n", o.EncryptionKeyFile)
      return false
   }

   if o.UseEncryption {
      if len(o.EncryptionKey) == 0 && len(o.EncryptionKeyFile) == 0 {
         fmt.Println("error: encryption key or encryption key file is required for encryption")
         return false
      }
   }

   return true
}
