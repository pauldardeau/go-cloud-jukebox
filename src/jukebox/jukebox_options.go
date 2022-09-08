package jukebox

import (
   "fmt"
)

type JukeboxOptions struct {
   Debug_mode bool
   Use_encryption bool
   Use_compression bool
   Check_data_integrity bool
   File_cache_count int
   Number_songs int
   Encryption_key string
   Encryption_key_file string
   Encryption_iv string
   Suppress_metadata_download bool
}

func NewJukeboxOptions() (*JukeboxOptions) {
   var o JukeboxOptions
   o.Debug_mode = false
   o.Use_encryption = false
   o.Use_compression = false
   o.Check_data_integrity = false
   o.File_cache_count = 3
   o.Number_songs = 0
   o.Encryption_key = ""
   o.Encryption_key_file = ""
   o.Encryption_iv = ""
   o.Suppress_metadata_download = false
   return &o
}

func (o *JukeboxOptions) Validate_options() (bool) {
   if o.File_cache_count < 0 {
      fmt.Println("error: file cache count must be non-negative integer value")
      return false
   }

   //TODO: add encryption support
   //if len(o.Encryption_key_file) > 0 && ! os.path.isfile(o.Encryption_key_file) {
   //   fmt.Printf("error: encryption key file doesn't exist '%s'\n", o.Encryption_key_file)
   //   return false
   //}

   //TODO: add encryption support
   //if o.Use_encryption {
   //   if ! aes.is_available() {
   //      fmt.Println("encryption support not available")
   //      return false
   //   }

   //   if len(o.Encryption_key) == 0 && len(o.Encryption_key_file) == 0 {
   //      fmt.Println("error: encryption key or encryption key file is required for encryption")
   //      return false
   //   }
   //}

   return true
}
