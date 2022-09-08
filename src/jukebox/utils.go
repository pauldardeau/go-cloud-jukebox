package jukebox

import (
   "crypto/md5"
   "errors"
   "fmt"
   "io"
   "os"
   "path/filepath"
)


func File_exists(path_to_file string) (bool) {
   if _, err := os.Stat(path_to_file); errors.Is(err, os.ErrNotExist) {
      return false
   } else {
      return true
   }
}

func Delete_file(path_to_file string) (bool) {
   err := os.Remove(path_to_file)
   return err == nil
}

func Path_join(dir_path string, file_name string) (string) {
   return filepath.Join(dir_path, file_name)
}

func Get_file_size(file_path string) int64 {
   fi, err := os.Stat(file_path)
   if err != nil {
      return -1
   }
   return fi.Size()
}

func Md5_for_file(path_to_file string) (string, error) {
   f, err := os.Open(path_to_file)
   if err != nil {
      return "", err
   }

   defer f.Close()
   h := md5.New()
   if _, err := io.Copy(h, f); err != nil {
      return "", err
   }
   return fmt.Sprintf("%x", h.Sum(nil)), nil
}

